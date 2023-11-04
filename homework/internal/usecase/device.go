package usecase

import (
	"homework/internal/domain"
)

type DeviceUseCase interface {
	GetDevice(string) (domain.Device, error)
	CreateDevice(d domain.Device) error
	DeleteDevice(string) error
	UpdateDevice(domain.Device) error
}
