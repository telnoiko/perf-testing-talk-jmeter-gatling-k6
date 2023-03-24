package store

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type Store struct {
	pool *pgxpool.Pool
}

func New() (*Store, error) {
	pool, err := pgxpool.New(context.Background(), "postgres://postgres:password@db:5432/tasks?pool_max_conns=10")
	if err != nil {
		log.Printf("Unable to create connection pool: %v\n", err)
		return nil, err
	}
	return &Store{pool}, nil
}

func (s *Store) Save(user *User) (int, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	var id int
	_, err = s.pool.Exec(context.Background(), "INSERT INTO users (name, email, password, tokens) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, hash, user.Tokens)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return 0, err
	}
	return id, nil
}

func (s *Store) UpdateToken(id int, token string) error {
	tag, err := s.pool.Exec(context.Background(), "UPDATE users SET tokens = array_append(tokens, $1) WHERE id = $2", token, id)
	if err != nil {
		log.Printf("Exec failed: %v\n", err)
		return err
	}
	return nil
}

func (s *Store) FindByEmail(email string) (*User, error) {
	var u User
	err := s.pool.QueryRow(context.Background(), "SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&u.ID, &u.Name, &u.Email, &u.Password)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		return &User{}, err
	}
	return &u, nil
}

func (s *Store) Stop() {
	s.pool.Close()
}
