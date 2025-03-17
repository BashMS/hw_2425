package storage

type User struct {
	ID      string
	Name    string
	Address string
}

type UserRepo interface {
	CreateUser(user User) (int64, error)
	UpdateUser(user User) error
	DeleteUser(userId int64) error
}