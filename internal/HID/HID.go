package HID

import (
	"github.com/karalabe/hid"
)

type ()

func GetDevices() []hid.DeviceInfo {
	return hid.Enumerate(1118, 0)
}
