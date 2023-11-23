package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"homework/internal/domain"
	"homework/internal/usecase"
	"net"
	"net/http"
)

// Handler TODO: определить набор полей и методов
type Handler struct {
	deviceUC usecase.DeviceUseCase
}

func NewHandler(deviceUC usecase.DeviceUseCase) *Handler {
	return &Handler{
		deviceUC: deviceUC,
	}
}
func (h *Handler) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var device domain.Device
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if net.ParseIP(device.IP).To4() == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.deviceUC.CreateDevice(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serialNum := params["serialNum"]

	device, err := h.deviceUC.GetDevice(serialNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(device)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serialNum := params["serialNum"]

	err := h.deviceUC.DeleteDevice(serialNum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) UpdateDevice(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	serialNum := params["serialNum"]

	var updatedDevice domain.Device
	err := json.NewDecoder(r.Body).Decode(&updatedDevice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedDevice.SerialNum = serialNum

	err = h.deviceUC.UpdateDevice(updatedDevice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RegisterHandlers(router *mux.Router) {
	router.HandleFunc("/api/v1/devices/{serialNum}", h.GetDevice).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/devices", h.CreateDevice).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/devices/{serialNum}", h.DeleteDevice).Methods(http.MethodDelete)
	router.HandleFunc("/api/v1/devices/{serialNum}", h.UpdateDevice).Methods(http.MethodPut)
}
