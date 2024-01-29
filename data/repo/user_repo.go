package repo

import (
	"time"

	"github.com/lengocson131002/go-clean/internal/domain"
	repo "github.com/lengocson131002/go-clean/internal/interfaces"
	"github.com/lengocson131002/go-clean/pkg/database"
)

type UserRepository struct {
	DB database.SqlGdbc
}

func NewUserRepository(tds database.TxDataInterface) *UserRepository {
	return &UserRepository{tds.GetTx()}
}

var _ repo.UserRepositoryInterface = (*UserRepository)(nil)

// CountById implements repo.UserRepositoryInterface.
func (r *UserRepository) CountById(id string) (int64, error) {
	var total int64
	err := r.DB.Get(&total, "SELECT * FROM users WHERE id = $1", id)
	return total, err
}

// CreateUser implements repo.UserRepositoryInterface.
func (r *UserRepository) CreateUser(user *domain.User) error {
	sql := "INSERT INTO users(id, password, name, created_at, updated_at) VALUES ($1, $2, $3, $4, $5)"
	_, err := r.DB.Exec(sql, user.ID, user.Password, user.Name, user.Token, time.Now().UTC().UnixMilli(), time.Now().UTC().UnixMilli())
	return err
}

// FindByToken implements repo.UserRepositoryInterface.
func (r *UserRepository) FindByToken(token string) (*domain.User, error) {
	// userEntity := new(entity.UserEntity)
	// err := r.db.Where("token = ?", token).First(userEntity).Error
	// return mapper.ToUserDomain(userEntity), err
	panic("failed")
}

// FindUserById implements repo.UserRepositoryInterface.
func (r *UserRepository) FindUserById(id string) (*domain.User, error) {
	// entity, err := r.FindById(id)
	// if err != nil {
	// 	return nil, err
	// }

	// return mapper.ToUserDomain(entity), err
	panic("failed")
}

// UpdateUser implements repo.UserRepositoryInterface.
func (r *UserRepository) UpdateUser(user *domain.User) error {
	// entity := mapper.ToUserEntity(user)
	// err := r.Update(entity)
	// return err
	panic("failed")
}
