package user

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	Get(id string) (*User, error)
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) error {
	user.ID = uuid.New().String()

	result := repo.db.Create(user)

	if result.Error != nil {
		repo.log.Println(result.Error)
		return result.Error
	}

	repo.log.Println("user created with id: ", user.ID)
	return nil
}

func (repo *repo) GetAll() ([]User, error) {
	var users []User

	result := repo.db.Model(&users).Order("created_at desc").Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	return users, nil
}

func (repo *repo) Get(id string) (*User, error) {
	user := User{ID: id}

	result := repo.db.First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}