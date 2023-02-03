package bl

import (
	"github.com/jkrus/master_api/internal"
)

type BL struct{}

func NewBL(di internal.IAppDeps) *BL {
	return &BL{}
}
