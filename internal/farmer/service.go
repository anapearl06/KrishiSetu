package farmer

import (
	"errors"
	"time"

	"github.com/raaj2493/KrishiSetu/internal/auth"
)

type Service struct {
	repo      Repository
	jwtSecret string
	jwtExpiry time.Duration
}

type RegisterInput struct {
	Name     string
	Phone    string
	Password string
	State    string
	District string
}

type LoginInput struct {
	Phone    string
	Password string
}

type LoginResult struct {
	Farmer *Farmer
	Token  string
}

func NewService(
	repo Repository,
	jwtSecret string,
	jwtExpiry time.Duration,
) *Service {
	return &Service{
		repo:      repo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
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

func (s *Service) Login(input LoginInput) (*LoginResult, error) {
	farmer, err := s.repo.FindByPhone(input.Phone)
	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(input.Password, farmer.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(
		farmer.ID,
		"farmer",
		s.jwtSecret,
		s.jwtExpiry,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Farmer: farmer,
		Token:  token,
	}, nil
}