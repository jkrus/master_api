package ping_controller

import "github.com/jkrus/master_api/internal/bl"

type PingController struct {
	bl *bl.BL
}

func NewPingController(bl *bl.BL) *PingController {
	return &PingController{bl: bl}
}
