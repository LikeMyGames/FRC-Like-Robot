module tennis-ball-shooter

go 1.24.0

require github.com/LikeMyGames/FRC-Like-Robot/state v0.0.0-20250924013705-647f9ffa42ef

replace github.com/LikeMyGames/FRC-Like-Robot/state => ./../state

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/warthog618/go-gpiocdev v0.9.1 // indirect
	golang.org/x/sys v0.18.0 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
)
