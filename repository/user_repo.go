package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lengocson131002/go-clean/domain"
	"github.com/lengocson131002/go-clean/pkg/database"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
)

type UserEntity struct {
	ID        string  `db:"id"`
	Password  string  `db:"password"`
	Name      string  `db:"name"`
	Token     *string `db:"token"`
	CreatedAt int64   `db:"created_at"`
	UpdatedAt int64   `db:"updated_at"`
}

type userRepositoryGdbc struct {
	DB *database.Gdbc
}

func NewUserRepository(g *database.Gdbc) *userRepositoryGdbc {
	return &userRepositoryGdbc{g}
}

// WithinTransaction implements repo.UserRepositoryInterface.
func (ur *userRepositoryGdbc) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return ur.DB.WithinTransaction(ctx, txFunc)
}

func (ur *userRepositoryGdbc) WithinTransactionOptions(ctx context.Context, txFunc func(ctx context.Context) error, txOptions *sql.TxOptions) error {
	return ur.DB.WithinTransactionOptions(ctx, txFunc, txOptions)
}

// FindByToken implements repo.UserRepositoryInterface.
func (ur *userRepositoryGdbc) FindByToken(ctx context.Context, token string) (*domain.User, error) {
	userEntity := &UserEntity{}
	err := ur.DB.Get(ctx, userEntity, "SELECT * FROM users WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	res := &domain.User{}
	err = mapper.BindingStruct(userEntity, &res)
	return res, err
}

// FindUserById implements repo.UserRepositoryInterface.
func (ur *userRepositoryGdbc) FindUserById(ctx context.Context, id string) (*domain.User, error) {
	userEntity := &UserEntity{}
	err := ur.DB.Get(ctx, userEntity, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	res := &domain.User{}
	err = mapper.BindingStruct(userEntity, &res)
	return res, err
}

// UpdateUser implements repo.UserRepositoryInterface.
func (ur *userRepositoryGdbc) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := ur.DB.Exec(ctx, "UPDATE users SET name=$1, password=$2, token=$3, updated_at=$4",
		user.Name,
		user.Password,
		user.Token,
		time.Now().UTC().UnixMilli())
	return err
}

var _ domain.UserRepository = (*userRepositoryGdbc)(nil)

// CountById implements repo.UserRepositoryInterface.
func (r *userRepositoryGdbc) CountById(ctx context.Context, id string) (int64, error) {
	var total int64
	err := r.DB.Get(ctx, &total, "SELECT COUNT(*) FROM users WHERE id = $1", id)
	return total, err
}

// CreateUser implements repo.UserRepositoryInterface.
func (r *userRepositoryGdbc) CreateUser(ctx context.Context, user *domain.User) error {
	sql := "INSERT INTO users(id, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.DB.Exec(ctx, sql,
		user.ID,
		user.Password,
		user.Name,
		time.Now().UTC().UnixMilli(),
		time.Now().UTC().UnixMilli())
	return err
}
