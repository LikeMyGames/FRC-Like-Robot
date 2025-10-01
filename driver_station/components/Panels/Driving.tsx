import style from "./Driving.module.css"
import { RobotRunSettingsContext } from "@/app/page"
import { useContext } from "react"

export default function DrivingPanel() {
    const [runSettings, setRunSettings] = useContext(RobotRunSettingsContext);

    return (
        <>
            <div className={style.runningmode_enabledisable_panel}>
                <div className={style.running_mode_selector}>
                    <button type="button" className={`${style.running_mode_selector_button} ${runSettings.mode == "teleop" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunSettings({ ...runSettings, mode: "teleop" })}>
                        TeleOperated
                    </button>
                    <button type="button" className={`${style.running_mode_selector_button} ${runSettings.mode == "auto" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunSettings({ ...runSettings, mode: "auto" })}>
                        Autonomous
                    </button>
                    <button type="button" className={`${style.running_mode_selector_button} ${runSettings.mode == "prac" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunSettings({ ...runSettings, mode: "prac" })}>
                        Practice
                    </button>
                    <button type="button" className={`${style.running_mode_selector_button} ${runSettings.mode == "test" ? style.running_mode_selector_button_selected : ""}`} onClick={() => setRunSettings({ ...runSettings, mode: "test" })}>
                        Test
                    </button>
                </div>
                <div className={style.enable_disable_panel}>
                    <button type="button" className={runSettings.enabled ? style.enable_disable_panel_button_selected : style.enable_disable_panel_button} onClick={() => setRunSettings({ ...runSettings, enabled: true })}>
                        Enable
                    </button>
                    <button type="button" className={!runSettings.enabled ? style.enable_disable_panel_button_selected : style.enable_disable_panel_button} onClick={() => setRunSettings({ ...runSettings, enabled: false })}>
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
        </>
    )
}