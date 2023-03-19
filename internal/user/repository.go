package user

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"strings"
)

type Repository interface {
	Create(user *User) error
	GetAll(filters Filters, offset, limit int) ([]User, error)
	Get(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName *string, lastName *string, email *string, phone *string) error
	Count(filters Filters) (int, error)
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

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]User, error) {
	var users []User

	tx := repo.db.Model(&users)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at desc").Find(&users)

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

func (repo *repo) Delete(id string) error {
	user := User{ID: id}

	result := repo.db.Delete(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repo) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {

	values := make(map[string]interface{})

	if firstName != nil {
		values["first_name"] = *firstName
	}

	if lastName != nil {
		values["last_name"] = *lastName
	}

	if email != nil {
		values["email"] = *email
	}

	if phone != nil {
		values["phone"] = *phone
	}

	result := repo.db.Model(&User{}).Where("id = ?", id).Updates(values)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(User{})
	tx = applyFilters(tx, filters)

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.FirstName != "" {
		filters.FirstName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.FirstName))
		tx = tx.Where("lower(trim(first_name)) like ?", filters.FirstName)
	}

	if filters.LastName != "" {
		filters.LastName = fmt.Sprintf("%%%s%%", strings.ToLower(filters.LastName))
		tx = tx.Where("lower(trim(last_name)) like ?", filters.LastName)
	}

	return tx
}
