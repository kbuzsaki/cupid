package server

type Server interface {
	KeepAlive() error
}
