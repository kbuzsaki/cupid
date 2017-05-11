package client

import "fmt"

type clientImpl struct {
	addr string
}

// TODO: accept a config file instead?
func New(addr string) (Client, error) {
	return &clientImpl{addr: addr}, nil
}

func (cl *clientImpl) Open(path string) (NodeHandle, error) {
	return nil, fmt.Errorf("Unable to open path: %#v", path)
}
