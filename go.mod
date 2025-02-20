module Robot

go 1.23.4

replace (
	github.com/LikeMyGames/FRC-Like-Robot/Controller => ./Controller
	github.com/LikeMyGames/FRC-Like-Robot/EventListener => ./EventListener
)

require (
	github.com/LikeMyGames/FRC-Like-Robot/Controller v0.0.0-00010101000000-000000000000
	github.com/LikeMyGames/FRC-Like-Robot/EventListener v0.0.0-00010101000000-000000000000
)

require github.com/tajtiattila/xinput v0.0.0-20140812192456-1b849e450040 // indirect
