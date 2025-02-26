module Robot

go 1.24.0

replace (
	internal/Controller => ./internal/Controller
	internal/EventListener => ./internal/EventListener
	internal/File => ./internal/File
	internal/Webpage => ./internal/Webpage
	internal/API => ./internal/API
)

require (
	internal/Controller v0.0.0
	internal/EventListener v0.0.0
	internal/Webpage v0.0.0
	internal/API v0.0.0
)

require (
	github.com/tajtiattila/xinput v0.0.0-20140812192456-1b849e450040 // indirect
	internal/File v0.0.0 // indirect
)
