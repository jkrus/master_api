package users_controller

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl"
)

type UserController struct {
	logger *zap.Logger
	bl     *bl.BL
}

func NewUserController(logger *zap.Logger, bl *bl.BL) *UserController {
	return &UserController{logger: logger, bl: bl}
}
