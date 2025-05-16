import style from "./RobotStatus.module.css"

export default function RobotStatus() {
    return (
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
    )
}