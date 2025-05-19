import { useContext } from "react"
import style from "./Settings.module.css"
import { useRobotContext } from "@/app/page"

export default function SettingsPanel() {
	const { RobotInfoContext } = useRobotContext()
	const [robotInfo, setRobotInfo] = useContext(RobotInfoContext)
	return (
		<div className={style.settings}>
			<div className={style.settings_panel}>
				<div className={style.settings_item}>
					<h3>Bot Net</h3>
					<input placeholder={"Bot Net"} defaultValue={robotInfo.botNet} onChange={(e) => { e.preventDefault(); setRobotInfo({ ...robotInfo, botNet: e.target.value }) }}></input>
				</div>
			</div>
			<div className={style.settings_panel}>
				<h3>Practice Timings (s)</h3>
				<div className={style.settings_item}>
					<h3>Countdown</h3>
					<input placeholder={"Team #"} defaultValue={1}></input>
				</div>
				<div className={style.settings_item}>
					<h3>Autonomous</h3>
					<input placeholder={"Team #"} defaultValue={1}></input>
				</div>
				<div className={style.settings_item}>
					<h3>Delay</h3>
					<input placeholder={"Team #"} defaultValue={1}></input>
				</div>
				<div className={style.settings_item}>
					<h3>Teleop</h3>
					<input placeholder={"Team #"} defaultValue={1}></input>
				</div>
				<div className={style.settings_item}>
					<h3>Endgame</h3>
					<input placeholder={"Team #"} defaultValue={1}></input>
				</div>
			</div>
		</div>
	)
}