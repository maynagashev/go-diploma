package repository

import (
	"github.com/jmoiron/sqlx"

	"gophermart/internal/domain"
)

// UserRepo реализует интерфейс domain.UserRepository.
type UserRepo struct {
	db *sqlx.DB
}

// NewUserRepo создает новый экземпляр UserRepo.
func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create добавляет нового пользователя в базу данных.
func (r *UserRepo) Create(user *domain.User) error {
	query := `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(
		query,
		user.Login,
		user.PasswordHash,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

// FindByLogin ищет пользователя по логину.
func (r *UserRepo) FindByLogin(login string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE login = $1`
	err := r.db.Get(&user, query, login)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
