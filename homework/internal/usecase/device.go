package usecase

import (
	"homework/internal/domain"
)

type DeviceUseCase interface {
	GetDevice(serialNum string) (d *domain.Device, err error)
	CreateDevice(d domain.Device) error
	DeleteDevice(serialNum string) error
	UpdateDevice(d domain.Device) error
}
