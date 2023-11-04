package repository

import (
	"fmt"
	"homework/internal/domain"
)

func (r *Repo) GetDevice(serialNum string) (d domain.Device, err error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.Devices[serialNum]
	if !ok {
		return domain.Device{}, fmt.Errorf("%w: no device", domain.ErrNotFound)
	}
	return d, nil
}

func (r *Repo) CreateDevice(d domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, e := r.Devices[d.SerialNum]
	if e {
		return fmt.Errorf("device is already in repository")
	}
	r.Devices[d.SerialNum] = d
	return nil
}
func (r *Repo) DeleteDevice(serialNum string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.Devices[serialNum]
	if !ok {
		return fmt.Errorf("%w: no device", domain.ErrNotFound)
	}
	delete(r.Devices, serialNum)
	return nil
}
func (r *Repo) UpdateDevice(d domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	_, e := r.Devices[d.SerialNum]
	if e {
		r.Devices[d.SerialNum] = d
		return nil
	}

	return fmt.Errorf("%w: no device", domain.ErrNotFound)
}
