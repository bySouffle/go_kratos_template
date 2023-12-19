package boot

import (
	"fmt"
	"go_kratos_template/app/template/internal/conf"
	"testing"
)

func TestWithID(t *testing.T) {
	b := BootStrap{
		Param: conf.Bootstrap{
			App: &conf.APPInfo{},
		},
	}
	b.Setting(WithID("10"))
	fmt.Println(b.Param.App)
}
