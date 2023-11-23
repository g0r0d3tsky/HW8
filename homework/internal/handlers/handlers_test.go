package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"homework/internal/domain"
	"homework/internal/handlers/mocks"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetDevice_TableDriven(t *testing.T) {

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
func TestCreateDevice_TableDriven(t *testing.T) {
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

func TestDeleteDevice_TableDriven(t *testing.T) {
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

func TestUpdateDevice_TableDriven(t *testing.T) {
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
			ExpectedStatus: http.StatusNoContent,
		},
		{
			Device: domain.Device{
				SerialNum: "2",
				Model:     "ppp",
				IP:        "0.9.9.0",
			},
			ExpectedStatus: http.StatusNoContent,
		},
	}

	for _, test := range testTable {
		mockDeviceUC.On("UpdateDevice", test.Device).Return(nil)

		deviceJSON, err := json.Marshal(test.Device)
		if err != nil {
			t.Errorf("error marshaling device: %s", err.Error())
		}

		req := httptest.NewRequest("PUT", "/devices/{serialNum}", bytes.NewBuffer(deviceJSON))
		recorder := httptest.NewRecorder()
		req = mux.SetURLVars(req, map[string]string{
			"serialNum": test.Device.SerialNum,
		})

		handler.UpdateDevice(recorder, req)

		if recorder.Code != test.ExpectedStatus {
			t.Errorf("expected status %d, got %d", test.ExpectedStatus, recorder.Code)
		}
	}
}
