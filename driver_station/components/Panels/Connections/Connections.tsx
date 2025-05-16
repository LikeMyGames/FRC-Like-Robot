import style from "./Connections.module.css"
import DeviceActionButtons from "./DeviceActionButtons"

export default function ConnectionsPanel() {
    return (
        <>
            <div className={style.device_panel}>
                <h3>USB Devices:</h3>
                <div className={style.device_container}>
                    <button className={style.device}>
                        <h3>Controller</h3>
                        <DeviceActionButtons />
                    </button>
                </div>
                <button className={style.device_rescan}>
                    <span className={"material-symbols-rounded"}>
                        refresh
                    </span>
                    Rescan
                </button>
            </div>
        </>
    )
}