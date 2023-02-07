package files_controller

import (
	"go.uber.org/zap"

	"github.com/jkrus/master_api/internal/bl"
)

type FileController struct {
	logger *zap.Logger
	bl     *bl.BL
}

func NewFileController(logger *zap.Logger, bl *bl.BL) *FileController {
	return &FileController{logger: logger, bl: bl}
}
