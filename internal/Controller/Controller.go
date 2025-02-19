package Controller

import (
	"strconv"

	"github.com/gookit/event"
	"github.com/tajtiattila/xinput"
)

func Main() {
	err := xinput.Load()
	if err == nil {
		var controllerState xinput.State
		var controller xinput.Gamepad
		var pressedButton []string
		for {
			xinput.GetState(0, &controllerState)
			controller = controllerState.Gamepad
			pressedButton = getPressedButton(controller.Buttons)
			for _, v := range pressedButton {
				event.Trigger(v, nil)
			}
		}
	}

}

func getPressedButton(sum uint16) []string {
	nums := []uint16{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 4096, 8192, 16384, 32768}
	str := []string{}
	for i := len(nums) - 1; i >= 0; i-- {
		if sum >= nums[i] {
			sum -= nums[i]
			str = append(str, strconv.Itoa(int(nums[i])))
		}
	}
	return str
}

func On(name string, callback func(e event.Event) error) {
	// event.On(name, event.ListenerFunc(func(e event.Event) error {
	// 	return nil
	// }), event.Normal)
}
