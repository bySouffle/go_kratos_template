package boot

import (
	"flag"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
)

var (
	flagConf string
)

type Conf struct {
	C config.Config
}

func NewBootConf() *Conf {
	return &Conf{}
}

func init() {
	// 获取命令行参数
	flag.StringVar(&flagConf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.Parse()
}

func (c *Conf) Load() {
	// 获取指定路径 yaml文件 的配置
	c.C = config.New(
		config.WithSource(
			file.NewSource(flagConf),
		),
	)

	if errLoad := c.C.Load(); errLoad != nil {
		panic(errLoad)
	}
}

func (c *Conf) Close() {
	c.C.Close()
}
