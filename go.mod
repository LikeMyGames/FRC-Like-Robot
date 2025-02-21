module main

go 1.23.4

replace (
	github.com/LikeMyGames/FRC-Like-Robot/internal/Controller => ./internal/Controller
	github.com/LikeMyGames/FRC-Like-Robot/internal/EventListener => ./internal/EventListener
	github.com/LikeMyGames/FRC-Like-Robot/internal/JSON => ./internal/JSON
)

require (
	github.com/LikeMyGames/FRC-Like-Robot/internal/Controller v0.0.0-00010101000000-000000000000
	github.com/LikeMyGames/FRC-Like-Robot/internal/EventListener v0.0.0-20250220072456-ca005ac1534d
)

require (
	github.com/LikeMyGames/FRC-Like-Robot/internal/JSON v0.0.0-20250220072456-ca005ac1534d // indirect
	github.com/tajtiattila/xinput v0.0.0-20140812192456-1b849e450040 // indirect
)
