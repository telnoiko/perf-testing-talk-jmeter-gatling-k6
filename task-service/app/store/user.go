package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	Tokens   []string `json:"tokens"`
}

type UserStore struct {
	pool   *pgxpool.Pool
	logger echo.Logger
}

func (s *UserStore) Create(user *User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = s.pool.Exec(context.Background(), "INSERT INTO users (name, email, password, tokens) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, hash, user.Tokens)
	if err != nil {
		s.logger.Printf("user create failed: %v\n", err)
		return err
	}
	return nil
}

func (s *UserStore) UpdateToken(id int, token string) error {
	_, err := s.pool.Exec(context.Background(), "UPDATE users SET tokens = array_append(tokens, $1) WHERE id = $2", token, id)
	if err != nil {
		s.logger.Printf("user token update failed: %v\n", err)
		return err
	}
	return nil
}

func (s *UserStore) FindByEmail(email string) (*User, error) {
	var u User
	err := s.pool.QueryRow(context.Background(), "SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		s.logger.Printf("find user by email failed: %v\n", err)
		return &User{}, err
	}
	return &u, nil
}

func (s *UserStore) FindByToken(email string, token string) (*User, error) {
	var u User
	row := s.pool.QueryRow(context.Background(), "SELECT id, name, email, password FROM users WHERE email = $1 AND tokens @> ARRAY[$2]::text[]", email, token)
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		s.logger.Printf("find user by email %d and token failed %s: %v\n", email, token, err)
		return &User{}, err
	}
	return &u, nil
}

func (s *UserStore) DeleteToken(id int, token string) error {
	_, err := s.pool.Exec(context.Background(), "UPDATE users SET tokens = array_remove(tokens, $1) WHERE id = $2", token, id)
	if err != nil {
		s.logger.Printf("delete user %d token failed: %v\n", id, err)
		return err
	}
	return nil
}
