package usecase

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework/internal/domain"
	"homework/internal/repository"
	"homework/internal/usecase/impl"
	"homework/internal/usecase/mocks"
	"testing"
)

func TestCreateDeviceMock(t *testing.T) {
	mockRepo := new(mocks.Device)
	useCase := &impl.UseCase{
		Repo: mockRepo,
	}

	mockRepo.On("CreateDevice", mock.Anything).Return(nil)
	err := useCase.CreateDevice(domain.Device{})
	mockRepo.AssertCalled(t, "CreateDevice", mock.Anything)
	assert.NoError(t, err)
}
func TestGetDevice(t *testing.T) {
	testCases := []struct {
		serialNum      string
		expectedDevice domain.Device
		expectedError  error
	}{
		{
			serialNum: "serialNumber1",
			expectedDevice: domain.Device{
				SerialNum: "serialNumber1",
				Model:     "xxx",
				IP:        "0.0.0.0",
			},
			expectedError: nil,
		},
		{
			serialNum: "serialNumber2",
			expectedDevice: domain.Device{
				SerialNum: "serialNumber2",
				Model:     "xxx",
				IP:        "0.0.0.0",
			},
			expectedError: nil,
		},
	}

	for _, tc := range testCases {
		mockRepo := new(mocks.Device)
		useCase := &impl.UseCase{
			Repo: mockRepo,
		}

		mockRepo.On("GetDevice", tc.serialNum).Return(tc.expectedDevice, tc.expectedError)

		device, err := useCase.GetDevice(tc.serialNum)

		mockRepo.AssertCalled(t, "GetDevice", tc.serialNum)

		assert.Equal(t, tc.expectedDevice, device)
		assert.Equal(t, tc.expectedError, err)
	}
}
func FuzzCreateDevice(f *testing.F) {
	repo := repository.New()
	service := impl.New(repo)
	f.Fuzz(func(t *testing.T, serialNum string) {
		if _, err := service.GetDevice(serialNum); err != nil {
			return
		}
		d := domain.Device{
			SerialNum: serialNum,
			Model:     "xxx",
			IP:        "0.0.0.0",
		}
		err := service.CreateDevice(d)
		if err != nil {
			t.Errorf("something wrong %v", err)
		}
	})

}
func TestCreateDevice(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	wantDevice := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}
	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotDevice, err := service.GetDevice(wantDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if wantDevice != gotDevice {
		t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
	}
}

func TestCreateMultipleDevices(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	devices := []domain.Device{
		{
			SerialNum: "123",
			Model:     "model1",
			IP:        "1.1.1.1",
		},
		{
			SerialNum: "124",
			Model:     "model2",
			IP:        "1.1.1.2",
		},
		{
			SerialNum: "125",
			Model:     "model3",
			IP:        "1.1.1.3",
		},
	}

	for _, d := range devices {
		err := service.CreateDevice(d)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	for _, wantDevice := range devices {
		gotDevice, err := service.GetDevice(wantDevice.SerialNum)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}

		if wantDevice != gotDevice {
			t.Errorf("want device %+#v not equal got %+#v", wantDevice, gotDevice)
		}
	}
}

func TestCreateDuplicate(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	wantDevice := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.CreateDevice(wantDevice)
	if err == nil {
		t.Errorf("want error, but got nil")
	}

}

func TestGetDeviceUnexisting(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	wantDevice := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(wantDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = service.GetDevice("1")
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDevice(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	newDevice := domain.Device{
		SerialNum: "123	",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	err = service.DeleteDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	_, err = service.GetDevice(newDevice.SerialNum)
	if err == nil {
		t.Error("want error, but got nil")
	}
}

func TestDeleteDeviceUnexisting(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)

	err := service.DeleteDevice("123")
	if err == nil {
		t.Errorf("want error, but got nil")
	}
}

func TestUpdateDevice(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	device := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(device)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	err = service.UpdateDevice(newDevice)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	gotDevice, err := service.GetDevice(newDevice.SerialNum)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if gotDevice != newDevice {
		t.Errorf("new device %+#v not equal got device %+#v", newDevice, gotDevice)
	}
}

func TestUpdateDeviceUnexsting(t *testing.T) {
	repo := repository.New()
	service := impl.New(repo)
	device := domain.Device{
		SerialNum: "123",
		Model:     "model1",
		IP:        "1.1.1.1",
	}

	err := service.CreateDevice(device)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	newDevice := domain.Device{
		SerialNum: "124",
		Model:     "model1",
		IP:        "1.1.1.2",
	}
	err = service.UpdateDevice(newDevice)
	if err == nil {
		t.Errorf("want err, but got nil")
	}
}
