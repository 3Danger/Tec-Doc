package web

type Server interface {
	Start() error
	Stop() error
}
