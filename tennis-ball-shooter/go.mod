module tennis-ball-shooter

go 1.24.0

require github.com/LikeMyGames/FRC-Like-Robot/state v0.0.0-20250924013705-647f9ffa42ef

replace github.com/LikeMyGames/FRC-Like-Robot/state => ./../state

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	periph.io/x/conn/v3 v3.7.2 // indirect
)
