package sd

type ServiceDiscovery interface {
	GetAddress(name string) string
	Register(onClose func(s *ServiceDiscovery)) error
	DeRegister() error
}


