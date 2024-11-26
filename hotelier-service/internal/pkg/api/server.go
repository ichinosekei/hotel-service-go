package api

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/ichinosekei/hotel-service-go/hotelier-service/internal/pkg/repository"
	"go.opentelemetry.io/otel/trace"
	"log"
	"net/http"
	"strconv"
)

// APIHandler  обработчик API, содержащий зависимости для передачи из main
type APIHandler struct {
	service *repository.Service
	tracer  trace.Tracer
}

// Создаёт новый APIHandler
func NewAPIHandler(service *repository.Service, tracer trace.Tracer) *APIHandler {
	return &APIHandler{service: service, tracer: tracer}
}

// Универсальная функция для обработки ошибок
func handleError(w http.ResponseWriter, statusCode int, err error) {
	log.Printf("Error: %v", err)
	http.Error(w, err.Error(), statusCode)
}

// Обработчики для работы с отелями

func (h *APIHandler) CreateHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "CreateHotelHandler")
	defer span.End()

	var data struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		handleError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	id, err := h.service.CreateHotel(data.Name, data.Location)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Hotel created: %d, name: %s, location: %s", id, data.Name, data.Location)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *APIHandler) GetHotelsHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "GetHotelsHandler")
	defer span.End()

	location := r.URL.Query().Get("location")
	hotels, err := h.service.GetHotels(location)
	if err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Retrieved %d hotels for location: %s", len(hotels), location)
	json.NewEncoder(w).Encode(hotels)
}

func (h *APIHandler) UpdateHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "UpdateHotelHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var data struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateHotel(id, data.Name, data.Location); err != nil {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	log.Printf("Hotel updated: %d, name: %s, location: %s", id, data.Name, data.Location)
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteHotelHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "DeleteHotelHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.DeleteHotel(id); err != nil {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}

	log.Printf("Hotel deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

// Обработчики для работы с комнатами

func (h *APIHandler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "CreateRoomHandler")
	defer span.End()

	var data struct {
		HotelID    int     `json:"hotel_id"`
		RoomNumber string  `json:"room_number"`
		Price      float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	//////////dsafsdfasfsdafsdf
	id, err := h.service.CreateRoom(data.HotelID, data.RoomNumber, data.Price)
	if err != nil {
		if err.Error() == "hotel not found" {
			handleError(w, http.StatusBadRequest, errors.New("hotel with the given ID does not exist"))
			return
		}
		handleError(w, http.StatusInternalServerError, err)
		return
	}

	log.Printf("Room created: %d, hotel_id: %d, room_number: %s, price: %.2f", id, data.HotelID, data.RoomNumber, data.Price)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *APIHandler) UpdateRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "UpdateRoomHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var data struct {
		RoomNumber string  `json:"room_number"`
		Price      float64 `json:"price"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateRoom(id, data.RoomNumber, data.Price); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	log.Printf("Room updated: %d, room_number: %s, price: %.2f", id, data.RoomNumber, data.Price)
	w.WriteHeader(http.StatusOK)
}

func (h *APIHandler) DeleteRoomHandler(w http.ResponseWriter, r *http.Request) {
	_, span := h.tracer.Start(r.Context(), "DeleteRoomHandler")
	defer span.End()

	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	if err := h.service.DeleteRoom(id); err != nil {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	log.Printf("Room deleted: %d", id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *APIHandler) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	// Начинаем span для трассировки
	_, span := h.tracer.Start(r.Context(), "GetRoomsHandler")
	defer span.End()

	// Получаем параметр `hotel_id` из строки запроса
	hotelIDParam := r.URL.Query().Get("hotel_id")
	var hotelID int
	var err error
	if hotelIDParam != "" {
		hotelID, err = strconv.Atoi(hotelIDParam)
		if err != nil {
			http.Error(w, "Invalid hotel_id parameter", http.StatusBadRequest)
			return
		}
	}

	// Получаем список комнат
	rooms, err := h.service.GetRooms(hotelID)
	if err != nil {
		http.Error(w, "Failed to retrieve rooms", http.StatusInternalServerError)
		return
	}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}
