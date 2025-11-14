'use client'
import style from "@/app/page.module.css";
import Icon from "@/components/Basic/Icon";
import Logger from "@/components/Main/Logger";
import RobotStatus from "@/components/Main/RobotStatus";
import ConnectionsPanel from "@/components/Panels/Connections/Connections";
import DrivingPanel from "@/components/Panels/Driving";
import SettingsPanel from "@/components/Panels/Settings";
import { useState, createContext, useRef, useEffect } from "react";
import { runExe } from "./scripts/controller-loader";

export type LoggerFilter = {
	all: boolean;
	log: boolean;
	success: boolean;
	warn: boolean;
	error: boolean;
	pause: boolean;
}

export type Log = {
	type: "log" | "success" | "warn" | "error" | "comment";
	message: string;
}

type SocketData = {
	system_logger?: Log;
	robot_status?: RobotStatus;
	run_settings?: RunSettings;
}

export type RobotStatus = {
	comms: boolean;
	code: boolean;
	joy: boolean;
	message: string;
	bat_p: number;
	bat_v: number;
}

export type RunSettings = {
	enabled: boolean;
	mode: string;
}

export type RobotInfo = {
	botNet: string
	controller: string
}

export const RobotRunSettingsContext = createContext<[RunSettings, (value: RunSettings) => void]>([{} as RunSettings, (value) => { console.log(value) }]);
export const LoggerContext = createContext<[Log[], () => void]>([[{ type: "success", message: "default message" }] as Log[], () => { }]);;
export const LoggerFilterContext = createContext<[LoggerFilter, (value: LoggerFilter) => void]>([{} as LoggerFilter, () => { }]);
export const RobotStatusContext = createContext<[RobotStatus]>([{} as RobotStatus]);
export const RobotInfoContext = createContext<[RobotInfo, (value: RobotInfo) => void]>([{} as RobotInfo, () => { }])

export default function Home() {
	const [panel, setPanel] = useState<string>("driving")
	const [robotInfo, setRobotInfo] = useState<RobotInfo>({ botNet: "ws://localhost:8080", controller: "ws://localhost:3030" } as RobotInfo)
	const robotConn = useRef<WebSocket | null>(null);
	const controllerConn = useRef<WebSocket | null>(null)
	const [runSettings, setRunSettings] = useState<RunSettings>({ mode: "teleop", enabled: false })
	const [loggerFilter, setLoggerFilter] = useState<LoggerFilter>({
		all: true,
		log: false,
		success: false,
		warn: false,
		error: false,
		pause: false,
		refresh: true,
	} as LoggerFilter)
	const [logs, setLogs] = useState<Log[]>([{ type: "warn", message: "if Robot Websocket connection is not attempted, press the logger window rehresh button" }] as Log[])
	const newLogs = useRef<Log[]>(logs)
	const [robotStat, setRobotStat] = useState<RobotStatus>({
		comms: false,
		code: false,
		joy: false,
		message: "TeleOperated Disabled",
		bat_p: 50,
		bat_v: 12.00
	} as RobotStatus)
	const newRobotStat = useRef<RobotStatus>(robotStat)
	const socketMessages = useRef<SocketData[]>([] as SocketData[])
	const lastSocketMessage = useRef<SocketData>({} as SocketData)

	const resetLogs = () => {
		newLogs.current = []
		setLogs([] as Log[])
	}

	function SetRobotRunSettings(settings: RunSettings) {
		if (settings.mode != runSettings.mode && runSettings.enabled) {
			settings.enabled = false
		}
		if (robotConn.current) {
			console.log("sending run_settings change")
			robotConn.current.send(`{"run_settings":${JSON.stringify(settings)}}`)
		}
		setRunSettings(settings)
	}

	useEffect(() => {
		function addLog(log: Log) {
			newLogs.current = [...newLogs.current, log]
			if (newLogs.current !== logs) {
				setLogs(newLogs.current)
			}
		}

		const setRobotStatus = (stat: RobotStatus) => {
			if (newRobotStat.current !== robotStat) {
				newRobotStat.current = { ...newRobotStat.current, ...stat } as RobotStatus
				setRobotStat(newRobotStat.current)
			} else {
				newRobotStat.current = { ...robotStat, ...stat } as RobotStatus
			}
		}

		if (controllerConn.current == null) {
			runExe("./ControllerInput/ControllerInput.exe")
			console.log("attempting controller re-connect")
			controllerConn.current = new WebSocket(robotInfo.controller)

			controllerConn.current.onmessage = (event) => {
				// console.log(event.data)
				if (event.data.system_logger) {
					addLog(event.data.system_logger as Log)
				}
				if (robotConn.current) {
					robotConn.current.send(event.data)
				}
			};

			controllerConn.current.onclose = () => {
				controllerConn.current = null
				addLog({ type: "warn", message: "controller disconnected" })
				setRobotStatus({ joy: false } as RobotStatus)
			};

			controllerConn.current.onerror = (error) => {
				addLog({ type: "error", message: "Controller Websocket connection error" })
				console.log('Controller WebSocket error:', error);
			};

			controllerConn.current.onopen = () => {
				setRobotStatus({ joy: true } as RobotStatus)
				addLog({ type: "success", message: "Controller WebSocket connected" })
				if (controllerConn.current) {
					controllerConn.current.send(`{"message":"controller connection good"}`)
				}
			};
		}

		if ((robotConn.current == null) && !loggerFilter.pause) {
			console.log("attempting websocket re-connect")
			robotConn.current = new WebSocket(robotInfo.botNet);

			robotConn.current.onmessage = (event) => {
				// Handle incoming messages
				const data = JSON.parse(event.data) as SocketData;
				lastSocketMessage.current = data
				socketMessages.current = [...socketMessages.current, lastSocketMessage.current]
				if (data.system_logger) {
					addLog(data.system_logger)
				}
				if (data.robot_status) {
					setRobotStatus(data.robot_status)
				}
				if (data.run_settings) {
					setRunSettings(data.run_settings)
				}
			};

			robotConn.current.onclose = () => {
				robotConn.current = null
				addLog({ type: "error", message: "Robot Websocket connection closed" })
				setRobotStatus({
					comms: false,
					code: false,
					message: "No Robot Connection",
					bat_p: 0,
					bat_v: 0
				} as RobotStatus)
				SetRobotRunSettings({ ...runSettings, enabled: false })
				console.log('WebSocket connection closed');
			};

			robotConn.current.onerror = (error) => {
				addLog({ type: "error", message: "Robot Websocket connection error" })
				console.log('WebSocket error:', error);
			};

			robotConn.current.onopen = () => {
				console.log('WebSocket connection established');
				addLog({ type: "success", message: "Robot Websocket connection established" })
				if (robotConn.current) {
					robotConn.current.send(`{"message":"GUI connected to robot"}`)
				}
				setRunSettings({ ...runSettings, enabled: false })
			};
		}

	}, [logs, robotInfo, robotConn, robotStat, newRobotStat, loggerFilter]);

	// useEffect(() => {
	// 	if (typeof window !== 'undefined') {
	// 		setIsWindowLoaded(true);
	// 	}
	// }, []);

	return (
		<div className={style.main_container}>
			<div className={style.sidemenu}>
				<button type="button" title="Driving" className={`${style.sidemenu_button} ${panel == "driving" ? style.button_selected : ""}`} onClick={() => setPanel("driving")}>
					<Icon iconName="search_hands_free" className={style.sidemenu_button_icon} />
				</button>
				<button type="button" title="Settings" className={`${style.sidemenu_button} ${panel == "settings" ? style.button_selected : ""}`} onClick={() => setPanel("settings")}>
					<Icon iconName="settings" className={style.sidemenu_button_icon} />
				</button>
				<button type="button" title="Connections" className={`${style.sidemenu_button} ${panel == "connections" ? style.button_selected : ""}`} onClick={() => setPanel("connections")}>
					<Icon iconName="usb" className={style.sidemenu_button_icon} />
				</button>
				<button type="button" title="Robot" className={`${style.sidemenu_button} ${panel == "robot" ? style.button_selected : ""}`} onClick={() => setPanel("robot")}>
					<Icon iconName="robot" className={style.sidemenu_button_icon} />
				</button>
			</div>
			<div className={style.panel_container}>
				{
					panel == "robot" ? (<>
						<RobotView />
					</>) : (<>
						<RobotInfoContext.Provider value={[robotInfo, setRobotInfo]}>
							<h1 className={style.panel_title}>{panel}</h1>
							<div className={style.panel_item_container}>
								<RobotRunSettingsContext.Provider value={[runSettings, SetRobotRunSettings]}>
									<div className={`${style.panel_item_subpanel} ${style.replaceable_panel}`}>
										{
											panel == "driving" ? (
												<DrivingPanel />
											) : panel == "settings" ? (
												<SettingsPanel />
											) : panel == "connections" ? (
												<ConnectionsPanel />
											) : (
												<>

												</>
											)
										}
									</div>
								</RobotRunSettingsContext.Provider>
								<RobotStatusContext.Provider value={[robotStat]}>
									<RobotStatus />
								</RobotStatusContext.Provider>
							</div>
							<LoggerContext.Provider value={[logs, resetLogs]}>
								<LoggerFilterContext.Provider value={[loggerFilter, setLoggerFilter]}>
									<Logger />
								</LoggerFilterContext.Provider>
							</LoggerContext.Provider>
						</RobotInfoContext.Provider></>)
				}
			</div>
		</div >
	)
}

export function useRobotContext() {
	return { RobotRunSettingsContext, RobotStatusContext, LoggerContext, LoggerFilterContext, RobotInfoContext }
}