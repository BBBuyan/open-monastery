package user

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetUser(ctx context.Context, id int64) (*User, error) {
	return s.repo.FindById(ctx,id)
}

func (s *Service) CreateUser(
	ctx context.Context,
	firstName string,
	lastName string,
) (*User, error) {
	user := &User{
		FirstName: firstName,
		LastName:  lastName,
		UserType:      Customer,
	}

	err := s.repo.Add(ctx, user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateMonk(
	ctx context.Context,
	firstName string,
	lastName string,
	monasteryId int64,
) (*User, error) {
	monk := &User{
		FirstName:   firstName,
		LastName:    lastName,
		UserType:        Monk,
		MonasteryId: &monasteryId,
	}

	err := s.repo.Add(ctx, monk)

	if err != nil {
		return nil, err
	}

	return monk, nil
}
