package repo

import (
	"context"
	"database/sql"
	"time"

	"github.com/lengocson131002/go-clean/data/entity"
	"github.com/lengocson131002/go-clean/internal/domain"
	"github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/pkg/database"
	mapper "github.com/lengocson131002/go-clean/pkg/util"
)

type UserRepository struct {
	DB *database.Gdbc
}

// WithinTransaction implements repo.UserRepositoryInterface.
func (ur *UserRepository) WithinTransaction(ctx context.Context, txFunc func(ctx context.Context) error) error {
	return ur.DB.WithinTransaction(ctx, txFunc)
}

func (ur *UserRepository) WithinTransactionOptions(ctx context.Context, txFunc func(ctx context.Context) error, txOptions *sql.TxOptions) error {
	return ur.DB.WithinTransactionOptions(ctx, txFunc, txOptions)
}

// FindByToken implements repo.UserRepositoryInterface.
func (ur *UserRepository) FindByToken(ctx context.Context, token string) (*domain.User, error) {
	userEntity := &entity.UserEntity{}
	err := ur.DB.Get(ctx, userEntity, "SELECT * FROM users WHERE token = $1", token)
	if err != nil {
		return nil, err
	}
	res := &domain.User{}
	err = mapper.BindingStruct(userEntity, &res)
	return res, err
}

// FindUserById implements repo.UserRepositoryInterface.
func (ur *UserRepository) FindUserById(ctx context.Context, id string) (*domain.User, error) {
	userEntity := &entity.UserEntity{}
	err := ur.DB.Get(ctx, userEntity, "SELECT * FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	res := &domain.User{}
	err = mapper.BindingStruct(userEntity, &res)
	return res, err
}

// UpdateUser implements repo.UserRepositoryInterface.
func (ur *UserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	_, err := ur.DB.Exec(ctx, "UPDATE users SET name=$1, password=$2, token=$3, updated_at=$4",
		user.Name,
		user.Password,
		user.Token,
		time.Now().UTC().UnixMilli())
	return err
}

func NewUserRepository(g *database.Gdbc) *UserRepository {
	return &UserRepository{g}
}

var _ interfaces.UserRepositoryInterface = (*UserRepository)(nil)

// CountById implements repo.UserRepositoryInterface.
func (r *UserRepository) CountById(ctx context.Context, id string) (int64, error) {
	var total int64
	err := r.DB.Get(ctx, &total, "SELECT COUNT(*) FROM users WHERE id = $1", id)
	return total, err
}

// CreateUser implements repo.UserRepositoryInterface.
func (r *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	sql := "INSERT INTO users(id, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.DB.Exec(ctx, sql,
		user.ID,
		user.Password,
		user.Name,
		time.Now().UTC().UnixMilli(),
		time.Now().UTC().UnixMilli())
	return err
}
