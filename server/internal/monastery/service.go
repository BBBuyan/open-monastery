package monastery

import "context"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetMonastery(ctx context.Context, id int64) (*Monastery, error) {
	return s.repo.FindById(ctx, id)
}

func (s *Service) CreateMonastery(
	ctx context.Context,
	name string,
) (*Monastery, error) {
	monastery := &Monastery{
		Name: name,
	}
	err := s.repo.Add(ctx, monastery)
	if err != nil {
		return nil, err
	}
	return monastery, nil
}
