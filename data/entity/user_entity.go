package entity

type UserEntity struct {
	ID        string `db:"id"`
	Password  string `db:"password"`
	Name      string `db:"name"`
	Token     string `db:"token"`
	CreatedAt int64  `db:"created_at"`
	UpdatedAt int64  `db:"updated_at"`
}

func (u *UserEntity) TableName() string {
	return "users"
}
