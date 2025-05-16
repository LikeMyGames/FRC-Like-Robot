'use client'
import style from "@/app/page.module.css";
import Icon from "@/components/Basic/Icon";
import Logger from "@/components/Main/Logger";
import RobotStatus from "@/components/Main/RobotStatus";
import ConnectionsPanel from "@/components/Panels/Connections/Connections";
import DrivingPanel from "@/components/Panels/Driving";
import SettingsPanel from "@/components/Panels/Settings";
import { useState, createContext, useRef, useEffect } from "react";

export type LoggerFilter = {
	all: boolean;
	log: boolean;
	success: boolean;
	warn: boolean;
	error: boolean;
}

export type Log = {
	type: "log" | "success" | "warn" | "error";
	message: string;
}

type SocketData = {
	system_logger?: Log;
	robot_status?: RobotStatus;
}

export type RobotStatus = {
	comms: boolean;
	code: boolean;
	joy: boolean;
	message: string;
	bat_p: number;
	bat_v: number;
}



export const RunningContext = createContext<[boolean, (value: boolean) => void]>([false, (value) => { console.log(value) }]);
export const RunningModeContext = createContext<[string, (value: string) => void]>(["teleop", (value) => { console.log(value) }]);
export const LoggerContext = createContext<[Log[], () => void]>([[{ type: "success", message: "default message" }] as Log[], () => { }]);;
export const LoggerFilterContext = createContext<[LoggerFilter, (value: LoggerFilter) => void]>([{} as LoggerFilter, () => { }]);
export const RobotStatusContext = createContext<[RobotStatus]>([{} as RobotStatus]);

export default function Home() {
	const [first, setFirst] = useState<boolean>(true)
	const [panel, setPanel] = useState<string>("driving")
	const [botNet] = useState<string>("ws://localhost:8080")
	const robotConn = useRef<WebSocket | null>(null);
	const [runningMode, setRunningMode] = useState<string>("teleop")
	const [running, setRunning] = useState<boolean>(false)
	const [loggerFilter, setLoggerFilter] = useState<LoggerFilter>({
		all: true,
		log: false,
		success: false,
		warn: false,
		error: false
	} as LoggerFilter)
	const [logs, setLogs] = useState<Log[]>([] as Log[])
	const newLogs = useRef<Log[]>(logs)
	const [robotStat, setRobotStat] = useState<RobotStatus>({} as RobotStatus)
	const socketMessages = useRef<SocketData[]>([] as SocketData[])
	const lastSocketMessage = useRef<SocketData>({} as SocketData)

	const resetLogs = () => {
		setLogs([] as Log[])
	}



	useEffect(() => {
		function addLog(log: Log) {
			newLogs.current = [...newLogs.current, log]
			if (newLogs.current !== logs) {
				setLogs(newLogs.current)
			}
		}

		if (first) {
			if (window) {
				window.addEventListener("beforeunload", () => {
					robotConn.current?.close()
				})
			}

			if (robotConn.current == null) {

				robotConn.current = new WebSocket(botNet);

				robotConn.current.onmessage = (event) => {
					// Handle incoming messages
					const data = JSON.parse(event.data) as SocketData;
					lastSocketMessage.current = data
					socketMessages.current = [...socketMessages.current, lastSocketMessage.current]
					if (data.system_logger != undefined) {
						addLog(data.system_logger)
					}
				};

				robotConn.current.onclose = () => {
					console.log('WebSocket connection closed');
				};

				robotConn.current.onerror = (error) => {
					console.error('WebSocket error:', error);
				};

				robotConn.current.onopen = () => {
					console.log('WebSocket connection established');
					if (robotConn.current) {
						robotConn.current.send(`{"message":"GUI connected to robot"}`)
					}
				};

				setFirst(false)
			}
		}
	}, [logs, botNet, robotConn, first]);

	return (
		<div className={style.main_container}>
			<div className={style.sidemenu}>
				<button type="button" title="Driving" className={`${style.sidemenu_button} ${panel == "driving" ? style.button_selected : ""}`} onClick={() => setPanel("driving")}
				>
					<Icon iconName="search_hands_free" className={style.sidemenu_button_icon} />
				</button>
				<button type="button" title="Settings" className={`${style.sidemenu_button} ${panel == "settings" ? style.button_selected : ""}`} onClick={() => setPanel("settings")}
				>
					<Icon iconName="settings" className={style.sidemenu_button_icon} />
				</button>
				<button type="button" title="Connections" className={`${style.sidemenu_button} ${panel == "connections" ? style.button_selected : ""}`} onClick={() => setPanel("connections")}
				>
					<Icon iconName="usb" className={style.sidemenu_button_icon} />
				</button>
			</div>
			<div className={style.panel_container}
			// selected="true"
			>
				<h1 className={style.panel_title}>{panel}</h1>
				<div className={style.panel_item_container}>
					<RunningContext.Provider value={[running, setRunning]}>
						<RunningModeContext.Provider value={[runningMode, setRunningMode]}>
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
						</RunningModeContext.Provider>
					</RunningContext.Provider>
					<RobotStatusContext.Provider value={[robotStat]}>
						<RobotStatus />
					</RobotStatusContext.Provider>
				</div>
				<LoggerContext.Provider value={[logs, resetLogs]}>
					<LoggerFilterContext.Provider value={[loggerFilter, setLoggerFilter]}>
						<Logger />
					</LoggerFilterContext.Provider>
				</LoggerContext.Provider>
			</div>
			{/* <div className="panel_container" title="Settings" selected="false">
                <h1 className="panel_title">Settings</h1>
                <div className="panel_item">

                </div>
                <div className="panel_item">

                </div>
                <div className="panel_item">

                </div>
            </div>
            <div className="panel_container" title="Connections" selected="false">
                <h1 className="panel_title">Connections</h1>
                <div className="panel_item">

                </div>
                <div className="panel_item">

                </div>
                <div className="panel_item">

                </div>
            </div> */}
		</div >
	)
}

export function useRobotContext() {
	return { RunningContext, RunningModeContext, RobotStatusContext, LoggerContext, LoggerFilterContext }
}