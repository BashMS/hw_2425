package internalhttp

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/app"     //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	"github.com/gorilla/mux"                                            //nolint:depguard
)

type MyHandler struct {
	App app.App
}

// srvResp структура ответа.
type srvResp struct {
	ID    int64  `json:"id,omitempty"`
	Error string `json:"error,omitempty"`
}

// response записывает ответ клиенту.
func response(w http.ResponseWriter, code int, resp interface{}) {
	wResp, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(code)
	w.Write(wResp)
}

// Hello возвращает сообщение Hello World.
func (h *MyHandler) Hello(w http.ResponseWriter, r *http.Request) { //nolint:revive
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

// CreateUser создает запись пользователя.
func (h *MyHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req storage.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	id, err := h.App.CreateUser(r.Context(), req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusCreated, srvResp{ID: id})
}

// UpdateUser изменяет информацию о пользователе.
func (h *MyHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req storage.User

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = h.App.UpdateUser(r.Context(), req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusOK, srvResp{ID: req.ID})
}

// DeleteUser удаляет запись пользователя.
func (h *MyHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["userID"]
	userID, err := strconv.ParseInt(params, 10, 2)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = h.App.DeleteUser(r.Context(), userID)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusOK, srvResp{})
}

// CreateEvent создает запись события.
func (h *MyHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req storage.Event

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	id, err := h.App.CreateEvent(r.Context(), req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusCreated, srvResp{ID: id})
}

// UpdateEvent изменяет информацию о событии.
func (h *MyHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	var req storage.Event

	body, err := io.ReadAll(r.Body)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = h.App.UpdateEvent(r.Context(), req)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusOK, srvResp{ID: req.ID})
}

// DeleteEvent удаляет событие.
func (h *MyHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)["eventID"]
	eventID, err := strconv.ParseInt(params, 10, 2)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	err = h.App.DeleteEvent(r.Context(), eventID)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}

	response(w, http.StatusOK, srvResp{})
}

// ListEventsForDay возвращает список событий за день.
func (h *MyHandler) ListEventsForDay(w http.ResponseWriter, r *http.Request) {
	var err error
	startDay := time.Now()
	dayStr := r.URL.Query().Get("startDay")
	if len(dayStr) != 0 {
		startDay, err = time.Parse("2.1.2006", dayStr)
		if err != nil {
			response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
			return
		}
	}

	resp, err := h.App.ListEventsForDay(r.Context(), startDay)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}
	if len(resp) == 0 {
		response(w, http.StatusNoContent, nil)
		return
	}

	response(w, http.StatusOK, resp)
}

// ListEventsForWeek возвращает список событий за неделю.
func (h *MyHandler) ListEventsForWeek(w http.ResponseWriter, r *http.Request) {
	var err error
	startDay := time.Now()
	dayStr := r.URL.Query().Get("startDay")
	if len(dayStr) != 0 {
		startDay, err = time.Parse("2.1.2006", dayStr)
		if err != nil {
			response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
			return
		}
	}

	resp, err := h.App.ListEventsForWeek(r.Context(), startDay)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}
	if len(resp) == 0 {
		response(w, http.StatusNoContent, nil)
		return
	}

	response(w, http.StatusOK, resp)
}

// ListEventsForMonth возвращает список событий за месяц.
func (h *MyHandler) ListEventsForMonth(w http.ResponseWriter, r *http.Request) {
	var err error
	startDay := time.Now()
	dayStr := r.URL.Query().Get("startDay")
	if len(dayStr) != 0 {
		startDay, err = time.Parse("2.1.2006", dayStr)
		if err != nil {
			response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
			return
		}
	}

	resp, err := h.App.ListEventsForMonth(r.Context(), startDay)
	if err != nil {
		response(w, http.StatusBadRequest, srvResp{Error: err.Error()})
		return
	}
	if len(resp) == 0 {
		response(w, http.StatusNoContent, nil)
		return
	}

	response(w, http.StatusOK, resp)
}
