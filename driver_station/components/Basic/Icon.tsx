// import style from "@/app/page.module.css";

export default function Icon({ iconName, className }: { iconName: string; className?: string; }) {
	return (
		<span className={`${className} material-symbols-rounded`}>
			{iconName}
		</span>
	)
}