package repository

import (
	"homework/internal/domain"
	"sync"
)

type Repo struct {
	Devices map[string]domain.Device
	mu      sync.RWMutex
}
type Device interface {
	GetDevice(string) (domain.Device, error)
	CreateDevice(d domain.Device) error
	DeleteDevice(string) error
	UpdateDevice(domain.Device) error
}

func New() *Repo {
	return &Repo{
		Devices: make(map[string]domain.Device),
	}
}
