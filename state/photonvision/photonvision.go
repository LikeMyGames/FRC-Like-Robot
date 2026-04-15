package photonvision

import (
	"fmt"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/levifitzpatrick1/go-nt4"
)

type (
	Camera struct {
		name        string
		cameraTopic *nt4.Topic
	}
)

var (
	cameras []*Camera
)

func NewCamera(name string) *Camera {
	c := new(Camera)
	c.name = name

	c.cameraTopic = conn.GetNT4Client().GetTopic(fmt.Sprintf("photonvision/%s", name))

	cameras = append(cameras, c)
	return c
}

func (c *Camera) ReadCameraTopic() {
	fmt.Println(c.cameraTopic.Properties)
}
