package repo

import (
	"github.com/lengocson131002/go-clean/data/entity"
	"github.com/lengocson131002/go-clean/data/mapper"
	"github.com/lengocson131002/go-clean/internal/domain"
	repo "github.com/lengocson131002/go-clean/internal/interfaces"

	"gorm.io/gorm"
)

type UserRepository struct {
	Repository[entity.UserEntity]
	db *gorm.DB
}

func NewuserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		Repository: Repository[entity.UserEntity]{db: db},
		db:         db,
	}
}

// CountById implements repo.UserRepositoryInterface.
func (r *UserRepository) CountById(id string) (int64, error) {
	var total int64
	err := r.db.Model(new(entity.UserEntity)).Where("id = ?", id).Count(&total).Error
	return total, err
}

// CreateUser implements repo.UserRepositoryInterface.
func (r *UserRepository) CreateUser(user *domain.User) error {
	userEntity := mapper.ToUserEntity(user)
	return r.Create(userEntity)
}

// FindByToken implements repo.UserRepositoryInterface.
func (r *UserRepository) FindByToken(token string) (*domain.User, error) {
	userEntity := new(entity.UserEntity)
	err := r.db.Where("token = ?", token).First(userEntity).Error
	return mapper.ToUserDomain(userEntity), err
}

// FindUserById implements repo.UserRepositoryInterface.
func (r *UserRepository) FindUserById(id string) (*domain.User, error) {
	entity, err := r.FindById(id)
	if err != nil {
		return nil, err
	}

	return mapper.ToUserDomain(entity), err
}

// UpdateUser implements repo.UserRepositoryInterface.
func (r *UserRepository) UpdateUser(user *domain.User) error {
	entity := mapper.ToUserEntity(user)
	err := r.Update(entity)
	return err
}

var _ repo.UserRepositoryInterface = (*UserRepository)(nil)
