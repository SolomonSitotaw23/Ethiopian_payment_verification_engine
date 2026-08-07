package providers

func init() {
	DefaultRegistry.Register(NewTelebirrProvider())
	DefaultRegistry.Register(NewCBEProvider())
	DefaultRegistry.Register(NewBOAProvider())
	DefaultRegistry.Register(NewAmharaProvider())
}
