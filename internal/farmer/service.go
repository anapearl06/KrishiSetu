package farmer

import "github.com/raaj2493/KrishiSetu/internal/auth"

type Service struct {
	repo Repository
}

type RegisterInput struct {
	Name     string
	Phone    string
	Password string
	State    string
	District string
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(input RegisterInput) (*Farmer, error) {
	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	f := &Farmer{
		Name:         input.Name,
		Phone:        input.Phone,
		PasswordHash: passwordHash,
		State:        input.State,
		District:     input.District,
	}

	if err := s.repo.Create(f); err != nil {
		return nil, err
	}

	return f, nil
}