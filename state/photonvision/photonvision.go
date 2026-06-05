package photonvision

import (
	"fmt"

	"github.com/LikeMyGames/FRC-Like-Robot/state/conn"
	"github.com/levifitzpatrick1/go-nt4"
)

type (
	Camera struct {
		name                       string
		path                       string
		cameraSubscription         *nt4.Subscription
		resultSubscriber           *nt4.Subscription
		resultList                 []*Result
		driverModePublisher        *nt4.Topic
		driverModeSubscriber       *nt4.Subscription
		driverMode                 bool
		fpsLimitPublisher          *nt4.Topic
		fpsLimitSubscriber         *nt4.Subscription
		fpsLimit                   float64
		heartbeatSubscriber        *nt4.Subscription
		heartbeat                  int
		cameraIntrinsicsSubscriber *nt4.Subscription
		cameraIntrinsics           []float64
		cameraDistortionSubsriber  *nt4.Subscription
		cameraDistortion           []float64
		topicName                  string
	}

	Result struct {
	}

	PoseEstimator struct {
		// fieldTags
		// tagModel
		// primaryStrategy
		// multiTagFallbackStrategy

	}
)

var (
	cameras []*Camera
)

func NewCamera(name string) *Camera {
	c := new(Camera)
	c.name = name
	c.path = fmt.Sprintf("/photonvision/%s", name)

	c.cameraSubscription = conn.NT4.Subscribe([]string{c.path}, nil)
	// conn.NT4.PublishBoolean(c.path + "/driverModeRequest", false)

	c.resultSubscriber = conn.NT4.Subscribe([]string{c.path + "/rawBytes"}, nil)

	c.driverModePublisher = conn.NT4.Publish(c.path+"/driverModeRequest", nt4.TypeBoolean, nil)
	c.driverModeSubscriber = conn.NT4.Subscribe([]string{c.path + "/driverMode"}, nil)
	c.driverModeSubscriber.SetCallback(func(topic *nt4.Topic, timestamp int64, value any) {
		c.driverMode = value.(bool)
	})

	c.fpsLimitPublisher = conn.NT4.Publish(c.path+"/fpsLimitRequest", nt4.TypeDouble, nil)
	c.fpsLimitSubscriber = conn.NT4.Subscribe([]string{c.path + "/fpsLimit"}, nil)
	c.fpsLimitSubscriber.SetCallback(func(topic *nt4.Topic, timestamp int64, value any) {
		c.fpsLimit = value.(float64)
	})

	c.heartbeatSubscriber = conn.NT4.Subscribe([]string{c.path + "/heartbeat"}, nil)
	c.heartbeatSubscriber.SetCallback(func(topic *nt4.Topic, timestamp int64, value any) {
		c.heartbeat = value.(int)
	})

	c.cameraIntrinsicsSubscriber = conn.NT4.Subscribe([]string{c.path + "/cameraIntrinsics"}, nil)
	c.cameraIntrinsicsSubscriber.SetCallback(func(topic *nt4.Topic, timestamp int64, value any) {
		c.cameraIntrinsics = value.([]float64)
	})
	c.cameraDistortionSubsriber = conn.NT4.Subscribe([]string{c.path + "/cameraDistortion"}, nil)
	c.cameraDistortionSubsriber.SetCallback(func(topic *nt4.Topic, timestamp int64, value any) {
		c.cameraDistortion = value.([]float64)
	})

	cameras = append(cameras, c)
	return c
}

func (c *Camera) GetDriverMode() bool {
	return c.driverMode
}

func (c *Camera) SetDriverMode(driverMode bool) {
	conn.NT4.SetValue(c.driverModePublisher, driverMode)
}

func (c *Camera) GetFpsLimit() float64 {
	return c.fpsLimit
}

func (c *Camera) SetFpsLimit(fpsLimit float64) {
	conn.NT4.SetValue(c.fpsLimitPublisher, fpsLimit)
}

func (c *Camera) GetAllUnreadResults() (results []Result) {
	for update := range c.resultSubscriber.Updates() {
		data, ok := (update.Value).([]byte)
		if !ok {
			conn.Log(fmt.Sprintf("Data from camera: %v", data))
		}
	}

	return results
}

// func (c *Camera) IsConnected() bool {
// 	curHeartbeat := c.heartbeat
// 	now := time.Now().UnixMilli()

// 	if curHeartbeat < 0 {
// 		// we have never heard from the camera
// 		return false
// 	}

// 	if curHeartbeat != prevHeartbeatValue {
// 		// New heartbeat value from the coprocessor
// 		prevHeartbeatChangeTime = now
// 		prevHeartbeatValue = curHeartbeat
// 	}

// 	// return (now - prevHeartbeatChangeTime) < HEARTBEAT_DEBOUNCE_SEC;
// }

func (c *Camera) Close() {
	c.cameraSubscription.CloseUpdates()
	c.cameraDistortionSubsriber.CloseUpdates()
	c.cameraIntrinsicsSubscriber.CloseUpdates()
	c.driverModeSubscriber.CloseUpdates()
	c.fpsLimitSubscriber.CloseUpdates()
	c.heartbeatSubscriber.CloseUpdates()
}
