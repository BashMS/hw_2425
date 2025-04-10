package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

func TestStorage_CreateUser(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	id, err := s.CreateUser(context.Background(), storage.User{
		ID:      0,
		Name:    "test",
		Address: "1111",
	})
	if err != nil {
		t.Errorf("Storage.CreateUser() error = %v", err)
	}
	if id != 1 {
		t.Errorf("Получили id %v; want 1", id)
	}

	id, err = s.CreateUser(context.Background(), storage.User{
		ID:      0,
		Name:    "test 2",
		Address: "1111",
	})
	if err == nil {
		t.Error("Ожидалась ошибка создания пользователя")
	}
	if id != 0 {
		t.Errorf("Получили id %v; want 0", id)
	}
}

func TestStorage_UpdateUser(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	s.Users[int64(1)] = storage.User{
		ID:      int64(1),
		Name:    "test",
		Address: "111",
	}

	err := s.UpdateUser(context.Background(), storage.User{
		ID:      1,
		Name:    "test",
		Address: "2222",
	})
	if err != nil {
		t.Errorf("Storage.UpdateUser() error = %v", err)
	}
	if s.Users[int64(1)].Address != "2222" {
		t.Errorf("Storage.UpdateUser() получили %v; ожидали 2222;", s.Users[int64(1)].Address)
	}

	err = s.UpdateUser(context.Background(), storage.User{
		ID:      0,
		Name:    "test",
		Address: "3333",
	})
	if err == nil {
		t.Error("Ожидалась ошибка обновления пользователя")
	}
}

func TestStorage_DeleteUser(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	s.Users[int64(1)] = storage.User{
		ID:      int64(1),
		Name:    "test",
		Address: "111",
	}

	err := s.DeleteUser(context.Background(), int64(1))
	if err != nil {
		t.Errorf("Storage.DeleteUser() error = %v", err)
	}
	if len(s.Users) > 0 {
		t.Errorf("Storage.DeleteUser(): ожидали пустой массив пользователей; %v", s.Users)
	}

	err = s.DeleteUser(context.Background(), int64(1))
	if err == nil {
		t.Error("Ожидалась ошибка удаления пользователя")
	}
}

func TestStorage_CreateEvent(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	startDate := time.Now()
	id, err := s.CreateEvent(context.Background(), storage.Event{
		ID:        int64(0),
		Name:      "test",
		StartDate: startDate,
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}
	if id != 1 {
		t.Errorf("Получили id %v; want 1", id)
	}

	id, err = s.CreateEvent(context.Background(), storage.Event{
		ID:        0,
		Name:      "test 2",
		StartDate: startDate,
		UserID:    int64(1),
	})
	if err == nil {
		t.Error("Ожидалась ошибка создания события")
	}
	if id != 0 {
		t.Errorf("Получили id %v; want 0", id)
	}
}

func TestStorage_UpdateEvent(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	startDate := time.Now()
	s.Events[int64(1)] = storage.Event{
		ID:        int64(1),
		Name:      "test",
		UserID:    int64(1),
		StartDate: startDate,
	}
	s.Events[int64(2)] = storage.Event{
		ID:        int64(1),
		Name:      "test 2",
		UserID:    int64(1),
		StartDate: startDate.Add(24 * time.Hour),
	}

	err := s.UpdateEvent(context.Background(), storage.Event{
		ID:        int64(1),
		Name:      "test 2",
		UserID:    int64(1),
		StartDate: startDate,
	})
	if err != nil {
		t.Errorf("Storage.UpdateEvent() error = %v", err)
	}
	if s.Events[int64(1)].Name != "test 2" {
		t.Errorf("Storage.UpdateEvent() получили %v; ожидали 'test 2';", s.Events[int64(1)].Name)
	}

	err = s.UpdateEvent(context.Background(), storage.Event{
		ID:        int64(2),
		Name:      "test 22",
		UserID:    int64(1),
		StartDate: startDate,
	})
	if err == nil {
		t.Errorf("Ожидалась ошибка обновления пользователя")
	}

	err = s.UpdateEvent(context.Background(), storage.Event{
		ID:        0,
		Name:      "test 3",
		UserID:    3,
		StartDate: time.Now(),
	})
	if err == nil {
		t.Error("Ожидалась ошибка обновления пользователя")
	}
}

func TestStorage_DeleteEvent(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	s.Events[int64(1)] = storage.Event{
		ID:     int64(1),
		Name:   "test",
		UserID: 1,
	}

	err := s.DeleteEvent(context.Background(), int64(1))
	if err != nil {
		t.Errorf("Storage.DeleteEvent() error = %v", err)
	}
	if len(s.Events) > 0 {
		t.Errorf("Storage.DeleteEvent(): ожидали пустой массив событий; %v", s.Events)
	}

	err = s.DeleteEvent(context.Background(), int64(1))
	if err == nil {
		t.Error("Ожидалась ошибка удаления события")
	}
}

func TestStorage_ListEventsForDay(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	startDate := time.Now()
	id, err := s.CreateEvent(context.Background(), storage.Event{
		ID:        int64(0),
		Name:      "test",
		StartDate: startDate,
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}
	if id != 1 {
		t.Errorf("Получили id %v; want 1", id)
	}

	_, err = s.CreateEvent(context.Background(), storage.Event{
		ID:        0,
		Name:      "test 2",
		StartDate: startDate.Add(24 * time.Hour),
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}

	list, err := s.ListEventsForDay(
		context.Background(),
		time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()),
	)
	if err != nil {
		t.Errorf("Storage.ListEventsForDay() error = %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Ожидалась одно событие, получили: %v", len(list))
	}
	if list[0].ID != 1 {
		t.Error("Ожидалась одно событие с id = 1")
	}
}

func TestStorage_ListEventsForWeek(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	startDate := time.Now()
	id, err := s.CreateEvent(context.Background(), storage.Event{
		ID:        int64(0),
		Name:      "test",
		StartDate: startDate,
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}
	if id != 1 {
		t.Errorf("Получили id %v; want 1", id)
	}

	_, err = s.CreateEvent(context.Background(), storage.Event{
		ID:        0,
		Name:      "test 2",
		StartDate: startDate.Add(7 * 24 * time.Hour),
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}

	list, err := s.ListEventsForWeek(
		context.Background(),
		time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()),
	)
	if err != nil {
		t.Errorf("Storage.ListEventsForWeek() error = %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Ожидалась одно событие, получили: %v", len(list))
	}
	if list[0].ID != 1 {
		t.Error("Ожидалась одно событие с id = 1")
	}
}

func TestStorage_ListEventsForMonth(t *testing.T) {
	s := New(config.Config{Source: "memory"}, logger.New("debug"))
	startDate := time.Now()
	id, err := s.CreateEvent(context.Background(), storage.Event{
		ID:        int64(0),
		Name:      "test",
		StartDate: startDate,
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}
	if id != 1 {
		t.Errorf("Получили id %v; want 1", id)
	}

	_, err = s.CreateEvent(context.Background(), storage.Event{
		ID:        0,
		Name:      "test 2",
		StartDate: startDate.Add(31 * 24 * time.Hour),
		UserID:    int64(1),
	})
	if err != nil {
		t.Errorf("Storage.CreateEvent() error = %v", err)
	}

	list, err := s.ListEventsForMonth(
		context.Background(),
		time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location()),
	)
	if err != nil {
		t.Errorf("Storage.ListEventsForMonth() error = %v", err)
	}

	if len(list) != 1 {
		t.Errorf("Ожидалась одно событие, получили: %v", len(list))
	}
	if list[0].ID != 1 {
		t.Error("Ожидалась одно событие с id = 1")
	}
}
