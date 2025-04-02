package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	_ "github.com/jackc/pgx/stdlib"                                     //nolint:depguard
)

// genNewID получает новый id из сиквенса БД.
func (s *Storage) genNewID(ctx context.Context) (int64, error) {
	var id int64
	query := "select nextval('public.sid') as 'id'"

	row := s.DB.QueryRowContext(ctx, query, nil)

	err := row.Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("row.Scan: %w", err)
	}

	return id, nil
}

// checkUser проверяет существование пользователя.
func (s *Storage) checkUser(ctx context.Context, user storage.User) (bool, error) {
	var cnt int64
	query := "select count(1) as cnt from public.tuser where lower(address) = lower($1)"

	row := s.DB.QueryRowContext(ctx, query, user.Address)

	err := row.Scan(&cnt)
	if err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}

	return cnt == 0, nil
}

// createUser создает пользователя.
func (s *Storage) createUser(ctx context.Context, user storage.User) (int64, error) {
	query := `insert into public.tuser(id, name, address)
	          values ($1, $2, $3)`

	_, err := s.DB.ExecContext(ctx, query, user.ID, user.Name, user.Address)
	if err != nil {
		return 0, fmt.Errorf("DB.ExecContext: %w", err)
	}

	return user.ID, nil
}

// updateUser обновляет данные пользователя.
func (s *Storage) updateUser(ctx context.Context, user storage.User) error {
	query := `update public.tuser
	             set name = $1, 
				     address = $2
	           where id = $3`

	_, err := s.DB.ExecContext(ctx, query, user.Name, user.Address, user.ID)
	if err != nil {
		return fmt.Errorf("DB.ExecContext: %w", err)
	}

	return nil
}

// deleteUser удаляет пользователя.
func (s *Storage) deleteUser(ctx context.Context, userID int64) error {
	query := `delete from public.tuser
	           where id = $1`

	_, err := s.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("DB.ExecContext: %w", err)
	}

	return nil
}

// checkEvent проверяет существование события в назначенное время.
func (s *Storage) checkEvent(ctx context.Context, event storage.Event) (bool, error) {
	var cnt int64
	query := `select count(1) as cnt
	            from public.tevent 
		       where id <> $1 
			     and user_id = $2 
				 and start_date_time = $3
	`
	row := s.DB.QueryRowContext(ctx, query, event.ID, event.UserID, event.StartDate)

	err := row.Scan(&cnt)
	if err != nil {
		return false, fmt.Errorf("row.Scan: %w", err)
	}

	return cnt == 0, nil
}

// createEvent создает событие.
func (s *Storage) createEvent(ctx context.Context, event storage.Event) (int64, error) {
	query := `insert into public.tevent(id, name, start_date_time, end_date_time, user_id, description, remind_for)
	          values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.DB.ExecContext(ctx,
		query,
		event.ID,
		event.Name,
		event.StartDate,
		event.EndDate,
		event.UserID,
		event.Description,
		event.RemindFor)
	if err != nil {
		return 0, fmt.Errorf("DB.ExecContext: %w", err)
	}

	return event.ID, nil
}

// updateEvent обновляет событие.
func (s *Storage) updateEvent(ctx context.Context, event storage.Event) error {
	query := `update public.tevent 
	             set name = $1, 
				     start_date_time = $2, 
				     end_date_time = $3, 
				     user_id = $4, 
				     description = $5,
				     remind_for = $6
			   where id = $7	 
	`

	_, err := s.DB.ExecContext(ctx,
		query,
		event.Name,
		event.StartDate,
		event.EndDate,
		event.UserID,
		event.Description,
		event.RemindFor,
		event.ID)
	if err != nil {
		return fmt.Errorf("DB.ExecContext: %w", err)
	}

	return nil
}

// deleteEvent удаляет событие.
func (s *Storage) deleteEvent(ctx context.Context, eventID int64) error {
	query := `delete from public.tevent
	           where id = $1`

	_, err := s.DB.ExecContext(ctx, query, eventID)
	if err != nil {
		return fmt.Errorf("DB.ExecContext: %w", err)
	}

	return nil
}

// getListEventByPeriod возвращает список событий за указанный период.
func (s *Storage) getListEventByPeriod(ctx context.Context, startDay, endDay time.Time) ([]storage.Event, error) {
	query := `select e.* 
	            from public.tevent e
			   where e.start_date_time >= :start
			     and e.end_date_time < :end
	`

	result := make([]storage.Event, 0)
	args := map[string]interface{}{
		"start": startDay,
		"end":   endDay,
	}
	rows, err := s.DB.NamedQueryContext(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		if err := rows.StructScan(&event); err != nil {
			return nil, fmt.Errorf("rows.StructScan: %w", err)
		}
		result = append(result, event)
	}

	return result, nil
}
