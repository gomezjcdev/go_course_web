package user

import "log"

type (
	Filters struct {
		FirstName string
		LastName  string
	}

	Service interface {
		Create(firstName, lastName, email, phone string) (*User, error)
		GetAll(filters Filters, offset, limit int) ([]User, error)
		Get(id string) (*User, error)
		Delete(id string) error
		Update(id string, firstName *string, lastName *string, email *string, phone *string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}
)

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

func (s service) GetAll(filters Filters, offset, limit int) ([]User, error) {
	users, err := s.repo.GetAll(filters, offset, limit)

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

func (s service) Count(filters Filters) (int, error) {
	return s.repo.Count(filters)
}

func (s service) Delete(id string) error {
	err := s.repo.Delete(id)

	if err != nil {
		return err
	}

	return nil
}

func (s service) Update(id string, firstName *string, lastName *string, email *string, phone *string) error {
	err := s.repo.Update(id, firstName, lastName, email, phone)

	if err != nil {
		return err
	}

	return nil
}
