import style from "./Connections.module.css"

export default function DeviceActionButtons() {
    return (
        <div className={style.device_action_buttons}>
            <button className={style.device_action_button}>
                <span className={"material-symbols-rounded"}>
                    keyboard_double_arrow_up
                </span>
            </button>
            <button className={style.device_action_button}>
                <span className={"material-symbols-rounded"}>
                    keyboard_arrow_up
                </span>
            </button>
            <button className={style.device_action_button}>
                <span className={"material-symbols-rounded"}>
                    keyboard_arrow_down
                </span>
            </button>
            <button className={style.device_action_button}>
                <span className={"material-symbols-rounded"}>
                    keyboard_double_arrow_down
                </span>
            </button>
        </div>
    )
}