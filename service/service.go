package service

type ServiceRunner interface {
	Start() error
	Stop() error
}
