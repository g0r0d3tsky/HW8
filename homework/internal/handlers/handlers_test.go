package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"homework/internal/domain"
	"homework/internal/handlers/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestHandler_GetDevice(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}
	type mockBehavior func(r *mocks.DeviceUseCase, device domain.Device)
	testTable := []struct {
		SerialNum            string
		ExpectedStatus       int
		ExpectedDevice       domain.Device
		mockBehavior         mockBehavior
		expectedResponseBody string
	}{
		{
			SerialNum:      "1",
			ExpectedStatus: http.StatusNotFound,
			ExpectedDevice: domain.Device{
				SerialNum: "1",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			mockBehavior: func(r *mocks.DeviceUseCase, expectedDevice domain.Device) {
				r.On("GetDevice", "1").Return(domain.Device{}, errors.New("can`t get device"))
			},
			expectedResponseBody: "can`t get device\n",
		},
		{
			SerialNum:      "2",
			ExpectedStatus: http.StatusOK,
			ExpectedDevice: domain.Device{
				SerialNum: "2",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			mockBehavior: func(r *mocks.DeviceUseCase, expectedDevice domain.Device) {
				r.On("GetDevice", "2").Return(expectedDevice, nil)
			},
			expectedResponseBody: "{\"SerialNum\":\"2\",\"Model\":\"ppp\",\"IP\":\"0.9.9.0\"}\n",
		},
	}

	for _, test := range testTable {
		test.mockBehavior(mockDeviceUC, test.ExpectedDevice)
		req := httptest.NewRequest("GET", "/devices/{serialNum}", nil)
		recorder := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"serialNum": test.SerialNum,
		})

		handler.GetDevice(recorder, req)

		if recorder.Code != test.ExpectedStatus {
			t.Errorf("expected status %d, got %d", test.ExpectedStatus, recorder.Code)
		}

		assert.Equal(t, recorder.Code, test.ExpectedStatus)
		assert.Equal(t, recorder.Body.String(), test.expectedResponseBody)
	}
}

func TestHandler_GetDeviceWithError(t *testing.T) {

	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}

	testTable := []struct {
		SerialNum      string
		ExpectedStatus int
		ExpectedDevice domain.Device
	}{
		{
			SerialNum:      "1",
			ExpectedStatus: http.StatusOK,
			ExpectedDevice: domain.Device{
				SerialNum: "1",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
		},
		{
			SerialNum:      "2",
			ExpectedStatus: http.StatusOK,
			ExpectedDevice: domain.Device{
				SerialNum: "2",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
		},
	}
	mockDeviceUC.On("GetDevice", "1").Return(testTable[0].ExpectedDevice, nil)
	mockDeviceUC.On("GetDevice", "2").Return(testTable[1].ExpectedDevice, nil)

	for _, test := range testTable {
		req := httptest.NewRequest("GET", "/devices/{serialNum}", nil)
		recorder := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"serialNum": test.SerialNum,
		})

		handler.GetDevice(recorder, req)

		if recorder.Code != test.ExpectedStatus {
			t.Errorf("expected status %d, got %d", test.ExpectedStatus, recorder.Code)
		}

		var device domain.Device
		err := json.Unmarshal(recorder.Body.Bytes(), &device)
		if err != nil {
			t.Errorf("error unmarshaling response body: %s", err.Error())
		}

		if !reflect.DeepEqual(device, test.ExpectedDevice) {
			t.Errorf("expected device %+v, got %+v", test.ExpectedDevice, device)
		}
	}
}
func TestHandler_CreateDevice(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}

	testTable := []struct {
		Device         domain.Device
		ExpectedStatus int
	}{
		{
			Device: domain.Device{
				SerialNum: "1",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			ExpectedStatus: http.StatusCreated,
		},
		{
			Device: domain.Device{
				SerialNum: "2",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			ExpectedStatus: http.StatusCreated,
		},
	}

	for _, test := range testTable {
		mockDeviceUC.On("CreateDevice", test.Device).Return(nil)

		deviceJSON, err := json.Marshal(test.Device)
		if err != nil {
			t.Errorf("error marshaling device: %s", err.Error())
		}

		req := httptest.NewRequest("POST", "/devices", bytes.NewBuffer(deviceJSON))
		recorder := httptest.NewRecorder()

		handler.CreateDevice(recorder, req)

		if recorder.Code != test.ExpectedStatus {
			t.Errorf("expected status %d, got %d", test.ExpectedStatus, recorder.Code)
		}
	}
}

func TestHandler_CreateDeviceWithUseCaseErr(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)

	handler := &Handler{
		deviceUC: mockDeviceUC,
	}

	device := domain.Device{
		SerialNum: "1",
		Model:     "ppp",
		IP:        "0.9.9.0",
	}

	expectedStatus := http.StatusConflict
	mockDeviceUC.On("CreateDevice", device).Return(errors.New("can`t create device"))

	deviceJSON, err := json.Marshal(device)
	if err != nil {
		t.Errorf("error marshaling device: %s", err.Error())
	}

	req := httptest.NewRequest("POST", "/devices", bytes.NewBuffer(deviceJSON))
	recorder := httptest.NewRecorder()
	handler.CreateDevice(recorder, req)

	if recorder.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, recorder.Code)
	}

}

func TestHandler_DeleteDevice(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}

	testTable := []struct {
		SerialNum      string
		ExpectedStatus int
	}{
		{
			SerialNum:      "1",
			ExpectedStatus: http.StatusNoContent,
		},
		{
			SerialNum:      "2",
			ExpectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range testTable {
		mockDeviceUC.On("DeleteDevice", test.SerialNum).Return(nil)

		req := httptest.NewRequest("DELETE", "/devices/{serialNum}", nil)
		recorder := httptest.NewRecorder()

		req = mux.SetURLVars(req, map[string]string{
			"serialNum": test.SerialNum,
		})

		handler.DeleteDevice(recorder, req)

		if recorder.Code != test.ExpectedStatus {
			t.Errorf("expected status %d, got %d", test.ExpectedStatus, recorder.Code)
		}
	}
}
func TestHandler_DeleteDeviceWithUCError(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}
	device := domain.Device{
		SerialNum: "1",
		Model:     "ppp",
		IP:        "0.9.9.0",
	}
	expectedStatus := http.StatusNotFound
	expectedError := fmt.Errorf("%w: no device", domain.ErrNotFound)
	mockDeviceUC.On("DeleteDevice", device.SerialNum).Return(expectedError)

	req := httptest.NewRequest("DELETE", "/devices/{serialNum}", nil)
	recorder := httptest.NewRecorder()

	req = mux.SetURLVars(req, map[string]string{
		"serialNum": device.SerialNum,
	})

	handler.DeleteDevice(recorder, req)

	assert.Equal(t, recorder.Code, expectedStatus)
}

func TestHandler_UpdateDevice(t *testing.T) {
	mockDeviceUC := new(mocks.DeviceUseCase)
	handler := &Handler{
		deviceUC: mockDeviceUC,
	}
	type mockBehavior func(r *mocks.DeviceUseCase, device domain.Device)
	testTable := []struct {
		Device               domain.Device
		ExpectedStatus       int
		mockBehavior         mockBehavior
		expectedResponseBody string
	}{
		{
			Device: domain.Device{
				SerialNum: "1",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			ExpectedStatus: http.StatusNotFound,
			mockBehavior: func(r *mocks.DeviceUseCase, device domain.Device) {
				r.On("UpdateDevice", device).Return(errors.New("can`t update device"))
			},
			expectedResponseBody: "can`t update device\n",
		},
		{
			Device: domain.Device{
				SerialNum: "2",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			ExpectedStatus: http.StatusNoContent,
			mockBehavior: func(r *mocks.DeviceUseCase, device domain.Device) {
				r.On("UpdateDevice", device).Return(nil)
			},
			expectedResponseBody: "",
		},
	}

	for _, test := range testTable {
		deviceJSON, err := json.Marshal(test.Device)
		test.mockBehavior(mockDeviceUC, test.Device)
		if err != nil {
			t.Errorf("error marshaling device: %s", err.Error())
		}

		req := httptest.NewRequest("PUT", "/devices/{serialNum}", bytes.NewBuffer(deviceJSON))
		recorder := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{
			"serialNum": test.Device.SerialNum,
		})

		handler.UpdateDevice(recorder, req)

		assert.Equal(t, recorder.Body.String(), test.expectedResponseBody)
		assert.Equal(t, recorder.Code, test.ExpectedStatus)
	}
}
