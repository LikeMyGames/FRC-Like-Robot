module Controller

go 1.23.4

replace (
	internal/EventListener => ../EventListener
	internal/File => ../File
	internal/Webpage => ../Webpage
)

require (
	github.com/tajtiattila/xinput v0.0.0-20140812192456-1b849e450040
	internal/EventListener v0.0.0-20250220072456-ca005ac1534d
	internal/File v0.0.0-20250220072456-ca005ac1534d
)
