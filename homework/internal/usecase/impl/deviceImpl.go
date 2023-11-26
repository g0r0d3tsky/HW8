package impl

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/repository"
)

type UseCase struct {
	Repo repository.Device
}

func (uc *UseCase) GetDevice(serialNum string) (*domain.Device, error) {
	device, err := uc.Repo.GetDevice(serialNum)
	if err != nil {
		return device, err
	}
	return device, nil
}

func (uc *UseCase) CreateDevice(d domain.Device) error {
	err := uc.Repo.CreateDevice(d)
	if err != nil {
		return fmt.Errorf("usecase createDevice %v", err)
	}
	return nil
}
func (uc *UseCase) DeleteDevice(serialNum string) error {
	err := uc.Repo.DeleteDevice(serialNum)
	if err != nil {
		return err
	}
	return nil
}
func (uc *UseCase) UpdateDevice(d domain.Device) error {
	err := uc.Repo.UpdateDevice(d)
	if err != nil {
		return err
	}
	return nil
}
func New(r *repository.Repo) *UseCase {
	return &UseCase{Repo: r}
}
