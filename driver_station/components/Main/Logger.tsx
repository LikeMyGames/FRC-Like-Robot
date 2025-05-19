import style from "./Logger.module.css"
import { useRobotContext, LoggerFilter, Log } from "@/app/page"
import { useContext, useEffect, useRef } from "react"
import Icon from "../Basic/Icon"



export default function Logger() {
	const { LoggerContext, LoggerFilterContext } = useRobotContext()
	const [logs, resetLogs] = useContext(LoggerContext)
	const displayedLogs = useRef<Log[]>(logs)
	const [loggerFilter, setLoggerFilter] = useContext(LoggerFilterContext);

	function calculateLoggerFilter(filter: LoggerFilter) {
		if (filter.log && filter.success && filter.warn && filter.error) {
			setLoggerFilter({
				all: true,
				log: false,
				success: false,
				warn: false,
				error: false
			} as LoggerFilter)
		} else {
			setLoggerFilter(filter)
		}
	}

	useEffect(() => {
		if (!loggerFilter.pause) {
			displayedLogs.current = logs
		}
	}, [logs, displayedLogs, loggerFilter])

	return (
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
							log: !loggerFilter.log
						} as LoggerFilter)
					}}>
						Log
					</button>
					<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.success && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="success" onClick={() => {
						calculateLoggerFilter({
							...loggerFilter,
							all: false,
							success: !loggerFilter.success
						} as LoggerFilter)
					}}>
						Success
					</button>
					<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.warn && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="warn" onClick={() => {
						calculateLoggerFilter({
							...loggerFilter,
							all: false,
							warn: !loggerFilter.warn
						} as LoggerFilter)
					}}>
						Warn
					</button>
					<button type="button" className={`${style.system_logger_filter_option} ${loggerFilter.error && !loggerFilter.all ? style.system_logger_filter_option_selected : ""}`} title="error" onClick={() => {
						calculateLoggerFilter({
							...loggerFilter,
							all: false,
							error: !loggerFilter.error
						} as LoggerFilter)
					}}>
						Error
					</button>
				</div>
				<div className={style.system_logger_action_container}>
					<button type="button" className={style.system_logger_filter_action} onClick={() => {
						if (loggerFilter.pause) {
							displayedLogs.current = logs
						}
						setLoggerFilter({ ...loggerFilter, pause: !loggerFilter.pause })
					}
					}>
						<Icon iconName="pause" />
						{/* <span className="material-symbols-rounded">
							pause
						</span> */}
					</button>
					<button type="button" className={style.system_logger_filter_action}
					// onClick="systemLogRefresh()"
					>
						<span className="material-symbols-rounded">
							refresh
						</span>
					</button>
					<button type="button" className={style.system_logger_filter_action} onClick={() => { resetLogs(); displayedLogs.current = [] }}>
						<span className="material-symbols-rounded">
							delete
						</span>
					</button>
				</div>
			</div>
			<div className={style.system_logger_display}>
				<div className={style.system_logger_display_container_text}>
					<div className={style.system_logger_display_text}>
						{
							displayedLogs.current.map((log: Log, i: number) => {
								switch (log.type) {
									case "log":
										if (loggerFilter.log || loggerFilter.all) {
											// console.log("adding ", log, " to system logger")
											return (<h3 key={i} className={`${style[log.type]}`}>LOG - {log.message}</h3>)
										}
									case "success":
										if (loggerFilter.success || loggerFilter.all) {
											// console.log("adding ", log, " to system logger")
											return (<h3 key={i} className={`${style[log.type]}`}>SUCCESS - {log.message}</h3>)
										}
									case "warn":
										if (loggerFilter.warn || loggerFilter.all) {
											// console.log("adding ", log, " to system logger")
											return (<h3 key={i} className={`${style[log.type]}`}>WARNING - {log.message}</h3>)
										}
									case "error":
										if (loggerFilter.error || loggerFilter.all) {
											// console.log("adding ", log, " to system logger")
											return (<h3 key={i} className={`${style[log.type]}`}>ERROR - {log.message}</h3>)
										}
									case "comment":
										return (<h3 key={i} className={`${style[log.type]}`}>{log.message}</h3>)
								}
							})
						}
					</div>
				</div>
			</div>
		</div>
	)
}