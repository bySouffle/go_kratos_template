package boot

type Boot interface {
	Load()
	Setting(opt ...interface{})
	Run()
	Close()
}
