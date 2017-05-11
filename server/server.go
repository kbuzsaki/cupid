package server

type serverImpl struct {
	addr string
}

// TODO: accept a config file instead?
func New(addr string) (Server, error) {
	return &serverImpl{addr: addr}, nil
}

func (s *serverImpl) KeepAlive() error {
	return nil
}
