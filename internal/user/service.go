package user

import "log"

type Service interface {
	Create(firstName, lastName, email, phone string) (*User, error)
	GetAll() ([]User, error)
	Get(id string) (*User, error)
}

type service struct {
	log  *log.Logger
	repo Repository
}

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(firstName, lastName, email, phone string) (*User, error) {
	s.log.Println("Create User Service")

	user := User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	err := s.repo.Create(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s service) GetAll() ([]User, error) {
	users, err := s.repo.GetAll()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s service) Get(id string) (*User, error) {
	user, err := s.repo.Get(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}