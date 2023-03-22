package store

type User struct {
	ID       int      `json:"-"`
	Name     string   `json:"name"`
	Email    string   `json:"email"`
	Password string   `json:"-"`
	Tokens   []string `json:"-"`
}
