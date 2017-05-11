package client

import "io"

type Client interface {
	Open(path string) (NodeHandle, error)
}

type Locker interface {
	Acquire() error
	TryAcquire() (bool, error)
	Release() error
}

type File interface {
	GetContentsAndStat() error
	GetStat() error
	SetContents(contents string) error
}

type Directory interface {
}

type NodeHandle interface {
	io.Closer
	Locker
	File
	Directory
	Delete() error
}
