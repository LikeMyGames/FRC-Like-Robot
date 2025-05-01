'use client'
import style from "@/app/page.module.css";
import Icon from "@/components/Basic/Icon";
import { useState } from "react";

type LoggerFilter = {
	all: boolean;
	log: boolean;
	success: boolean;
	warn: boolean;
	error: boolean;
}

export default function Home() {
	const [panel, setPanel] = useState<string>("driving")
	const [runningMode, setRunningMode] = useState<string>("teleop")
	const [running, setRunning] = useState<boolean>(false)
	const [loggerFilter, setLoggerFilter] = useState<LoggerFilter>({
		all: true,
		log: false,
		success: false,
		warn: false,
		error: false
	} as LoggerFilter)

	function calculateLoggerFilter(filter: LoggerFilter) {
		setLoggerFilter(filter)
	}

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
			<div className={style.panel_container} title="Driving"
			// selected="true"
			>
				<h1 className={style.panel_title}>{panel}</h1>
				<div className={style.panel_item_container}>
					<div className={`${style.panel_item_subpanel} ${style.driving_replaceable_panel}`}>
						<div className={style.runningmode_enabledisable_panel}>
							<div className={style.running_mode_selector}>
								<button type="button" className={`${style.running_mode_selector_button} ${runningMode == "teleop" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunningMode("teleop")}>
									TeleOperated
								</button>
								<button type="button" className={`${style.running_mode_selector_button} ${runningMode == "auto" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunningMode("auto")}>
									Autonomous
								</button>
								<button type="button" className={`${style.running_mode_selector_button} ${runningMode == "prac" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunningMode("prac")}>
									Practice
								</button>
								<button type="button" className={`${style.running_mode_selector_button} ${runningMode == "test" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunningMode("test")}>
									Test
								</button>
							</div>
							<div className={style.enable_disable_panel}>
								<button type="button" className={running ? style.enable_disable_panel_button_selected : style.enable_disable_panel_button} onClick={() => setRunning(true)}>
									Enable
								</button>
								<button type="button" className={!running ? style.enable_disable_panel_button_selected : style.enable_disable_panel_button} onClick={() => setRunning(false)}>
									Disable
								</button>
							</div>
						</div>
						<div className={style.elapsedtime_driverpos_panel}>
							<div className={style.elapsed_time_panel}>
								<h2>
									Elapsed Time:
								</h2>
								<h3 className={style.elapsed_time_value}>00:00</h3>
							</div>
							<div className={style.driver_pos_panel}>
								<h2>
									Team Station:
								</h2>
								<select className={style.driver_pos_select} title="driver_pos_select">
									<option className={style.driver_pos_option}>Red 1</option>
									<option className={style.driver_pos_option}>Red 2</option>
									<option className={style.driver_pos_option}>Red 3</option>
									<option className={style.driver_pos_option}>Blue 1</option>
									<option className={style.driver_pos_option}>Blue 2</option>
									<option className={style.driver_pos_option}>Blue 3</option>
								</select>
							</div>
						</div>
					</div>
					<div className={`${style.panel_item_subpanel} ${style.robot_status}`}>
						<div className={style.robot_battery_status}>
							<div className={style.robot_battery_status_display}>
								<div className={style.robot_battery_status_display_bar_container}>
									<div className={style.robot_battery_status_display_bar}></div>
									<h3 className={style.robot_battery_status_display_value}>50%</h3>
								</div>
							</div>
							<div className={style.robot_battery_status_voltage}>
								<h3 className={style.robot_battery_status_voltage_value}>12.00</h3>V
							</div>
						</div>
						<div className={style.robot_status_indicators}>
							<div className={style.robot_status_indicators_container}>
								<div className={style.robot_status_indicator}>
									<h3 className={style.robot_status_indicator_field}>Communications</h3>
									<div className={style.robot_status_indicator_display} />
								</div>
								<div className={style.robot_status_indicator}>
									<h3 className={style.robot_status_indicator_field}>Robot Code</h3>
									<div className={style.robot_status_indicator_display} />
								</div>
								<div className={style.robot_status_indicator}>
									<h3 className={style.robot_status_indicator_field}>Joysticks</h3>
									<div className={style.robot_status_indicator_display} />
								</div>
							</div>
						</div>
						<div className={style.robot_status_text}>
							TeleOperated Disabled
						</div>
					</div>
				</div>
				<div className={`${style.panel_item_container} ${style.system_logger}`}>
					<div className={style.system_logger_topmenu}>
						<div className={style.system_logger_filter_container}>
							<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="all" onClick={() => {
								calculateLoggerFilter({
									all: true,
									log: false,
									success: false,
									warn: false,
									error: false
								} as LoggerFilter)
							}}>
								All
							</button>
							<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.log && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="log" onClick={() => {
								calculateLoggerFilter({
									...loggerFilter,
									all: false,
									log: true
								} as LoggerFilter)
							}}>
								Log
							</button>
							<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.success && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="success" onClick={() => {
								calculateLoggerFilter({
									...loggerFilter,
									all: false,
									success: true
								} as LoggerFilter)
							}}>
								Success
							</button>
							<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.warn && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="warn" onClick={() => {
								calculateLoggerFilter({
									...loggerFilter,
									all: false,
									warn: true
								} as LoggerFilter)
							}}>
								Warn
							</button>
							<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.error && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="error" onClick={() => {
								calculateLoggerFilter({
									...loggerFilter,
									all: false,
									warn: true
								} as LoggerFilter)
							}}>
								Error
							</button>
						</div>
						<div className={style.system_logger_action_container}>
							<button type="button" className={style.system_logger_filter_action}
							// onClick="systemLogPause()"
							>
								<span className="material-symbols-rounded">
									pause
								</span>
							</button>
							<button type="button" className={style.system_logger_filter_action}
							// onClick="systemLogRefresh()"
							>
								<span className="material-symbols-rounded">
									refresh
								</span>
							</button>
							<button type="button" className={style.system_logger_filter_action}
							// onClick="systemLogClear()"
							>
								<span className="material-symbols-rounded">
									delete
								</span>
							</button>
						</div>
					</div>
					<div className={style.system_logger_display}>
						<div className={style.system_logger_display_text}>
							{/* <h4 className="log">LOG - Hello World</h4>
                            <h4 className="comment">COMMENT - Hello World</h4>
                            <h4 className="success">SUCCESS - Hello World</h4>
                            <h4 className="warn">WARNING - Hello World</h4>
                            <h4 className="error">ERROR - Hello World</h4> */}
						</div>
					</div>
				</div>
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