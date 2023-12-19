package boot

import (
	"github.com/go-kratos/kratos/v2/config"
	"go_kratos_template/app/auth/internal/conf"
)

type FolderOption func(*conf.Folder)

type Folder struct {
	C     config.Config
	Param conf.Folder
}

func NewFolder(c config.Config) *Folder {
	return &Folder{C: c}
}

func (d *Folder) Load() {
	if err := d.C.Scan(&d.Param); err != nil {
		panic(err)
	}
}

func (d *Folder) Setting(opts ...FolderOption) *Folder {
	for _, opt := range opts {
		opt(&d.Param)
	}
	return d
}
