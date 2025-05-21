package GuiControllerConnection

import (
	"frcrobot/internal/Controller"
)

func GetControllers() []*Controller.Controller {
	return Controller.Controllers
}
