package _interface

type Server interface {
	Start()
	GetNotify() <-chan error
	Shutdown()
}
