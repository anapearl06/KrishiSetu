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

type UpdateProfileInput struct {
	Name     string `json:"name"`
	State    string `json:"state"`
	District string `json:"district"`
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
	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.Phone == "" {
		return nil, errors.New("phone is required")
	}

	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	if input.State == "" {
		return nil, errors.New("state is required")
	}

	if input.District == "" {
		return nil, errors.New("district is required")
	}

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

func (s *Service) GetProfile(id uint) (*Farmer, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateProfile(
	id uint,
	input UpdateProfileInput,
) (*Farmer, error) {

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.State == "" {
		return nil, errors.New("state is required")
	}

	if input.District == "" {
		return nil, errors.New("district is required")
	}

	farmer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	farmer.Name = input.Name
	farmer.State = input.State
	farmer.District = input.District

	if err := s.repo.Update(farmer); err != nil {
		return nil, err
	}

	return farmer, nil
}

