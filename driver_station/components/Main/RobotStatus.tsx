import { useContext } from "react"
import style from "./RobotStatus.module.css"
import { useRobotContext } from "@/app/page"

export default function RobotStatus() {
    const { RobotStatusContext } = useRobotContext()
    const [robotStat] = useContext(RobotStatusContext)

    return (
        <div className={`${style.panel_item_subpanel} ${style.robot_status}`}>
            <div className={style.robot_battery_status}>
                <div className={style.robot_battery_status_display}>
                    <div className={style.robot_battery_status_display_bar_container}>
                        <div className={style.robot_battery_status_display_bar} style={{ width: `${robotStat.bat_p}%` }}></div>
                        <h3 className={style.robot_battery_status_display_value}>{robotStat.bat_p}%</h3>
                    </div>
                </div>
                <div className={style.robot_battery_status_voltage}>
                    <h3 className={style.robot_battery_status_voltage_value}>{robotStat.bat_v}</h3>V
                </div>
            </div>
            <div className={style.robot_status_indicators}>
                <div className={style.robot_status_indicators_container}>
                    <div className={style.robot_status_indicator}>
                        <h3 className={style.robot_status_indicator_field}>Communications</h3>
                        <div className={`${style.robot_status_indicator_display} ${robotStat.comms ? style.robot_status_indicator_display_true : ""}`} />
                    </div>
                    <div className={style.robot_status_indicator}>
                        <h3 className={style.robot_status_indicator_field}>Robot Code</h3>
                        <div className={`${style.robot_status_indicator_display} ${robotStat.code ? style.robot_status_indicator_display_true : ""}`} />
                    </div>
                    <div className={style.robot_status_indicator}>
                        <h3 className={style.robot_status_indicator_field}>Joysticks</h3>
                        <div className={`${style.robot_status_indicator_display} ${robotStat.joy ? style.robot_status_indicator_display_true : ""}`} />
                    </div>
                </div>
            </div>
            <div className={style.robot_status_text}>
                {robotStat.message}
            </div>
        </div>
    )
}