import style from "./Driving.module.css"
import { RunningContext, RunningModeContext } from "@/app/page"
import { useContext } from "react"

export default function DrivingPanel() {
    const [runningMode, setRunningMode] = useContext(RunningModeContext);
    const [running, setRunning] = useContext(RunningContext);

    return (
        <>
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
        </>
    )
}