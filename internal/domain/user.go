package domain

type User struct {
	ID        string
	Password  string
	Name      string
	Token     string
	CreatedAt int64
	UpdatedAt int64
}
