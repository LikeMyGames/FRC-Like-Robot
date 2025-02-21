module Controller

go 1.23.4

replace (
	github.com/LikeMyGames/FRC-Like-Robot/EventListener => ../EventListener
	github.com/LikeMyGames/FRC-Like-Robot/JSON => ../JSON
)

require (
	github.com/LikeMyGames/FRC-Like-Robot/internal/EventListener v0.0.0-20250220072456-ca005ac1534d
	github.com/LikeMyGames/FRC-Like-Robot/internal/JSON v0.0.0-20250220072456-ca005ac1534d
	github.com/tajtiattila/xinput v0.0.0-20140812192456-1b849e450040
)
