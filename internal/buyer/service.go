package buyer

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/raaj2493/KrishiSetu/internal/auth"
)

type Service struct {
	repo      Repository
	jwtSecret string
	jwtExpiry time.Duration
}

type RegisterInput struct {
	Name         string `json:"name"`
	Phone        string `json:"phone"`
	Password     string `json:"password"`
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`
	State        string `json:"state"`
	District     string `json:"district"`
}
type LoginInput struct {
	Phone    string
	Password string
}

type UpdateProfileInput struct {
	Name         string `json:"name"`
	BusinessName string `json:"business_name"`
	BusinessType string `json:"business_type"`
	State        string `json:"state"`
	District     string `json:"district"`
}

type LoginResult struct {
	Buyer *Buyer
	Token string
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

func (s *Service) Register(ctx context.Context, input RegisterInput) (*Buyer, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Phone = strings.TrimSpace(input.Phone)
	input.BusinessName = strings.TrimSpace(input.BusinessName)
	input.BusinessType = strings.TrimSpace(input.BusinessType)
	input.State = strings.TrimSpace(input.State)
	input.District = strings.TrimSpace(input.District)

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.Phone == "" {
		return nil, errors.New("phone is required")
	}

	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	if input.BusinessName == "" {
		return nil, errors.New("business name is required")
	}

	if input.BusinessType == "" {
		return nil, errors.New("business type is required")
	}

	if input.State == "" {
		return nil, errors.New("state is required")
	}

	if input.District == "" {
		return nil, errors.New("district is required")
	}

	_, err := s.repo.FindByPhone(ctx, input.Phone)
	if err == nil {
		return nil, errors.New("phone number already registered")
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	b := &Buyer{
		Name:         input.Name,
		Phone:        input.Phone,
		PasswordHash: passwordHash,
		BusinessName: input.BusinessName,
		BusinessType: input.BusinessType,
		State:        input.State,
		District:     input.District,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *Service) Login(ctx context.Context, input LoginInput) (*LoginResult, error) {
	input.Phone = strings.TrimSpace(input.Phone)

	if input.Phone == "" {
		return nil, errors.New("phone is required")
	}

	if input.Password == "" {
		return nil, errors.New("password is required")
	}

	buyer, err := s.repo.FindByPhone(ctx, input.Phone)
	if err != nil {
		return nil, err
	}

	if !auth.CheckPassword(input.Password, buyer.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	token, err := auth.GenerateToken(
		buyer.ID,
		"buyer",
		s.jwtSecret,
		s.jwtExpiry,
	)
	if err != nil {
		return nil, err
	}

	return &LoginResult{
		Buyer: buyer,
		Token: token,
	}, nil
}

func (s *Service) GetProfile(ctx context.Context, id uint) (*Buyer, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *Service) UpdateProfile(
	ctx context.Context,
	id uint,
	input UpdateProfileInput,
) (*Buyer, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.BusinessName = strings.TrimSpace(input.BusinessName)
	input.BusinessType = strings.TrimSpace(input.BusinessType)
	input.State = strings.TrimSpace(input.State)
	input.District = strings.TrimSpace(input.District)

	if input.Name == "" {
		return nil, errors.New("name is required")
	}

	if input.BusinessName == "" {
		return nil, errors.New("business name is required")
	}

	if input.State == "" {
		return nil, errors.New("state is required")
	}

	if input.District == "" {
		return nil, errors.New("district is required")
	}

	buyer, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	buyer.Name = input.Name
	buyer.BusinessName = input.BusinessName
	buyer.BusinessType = input.BusinessType
	buyer.State = input.State
	buyer.District = input.District

	if err := s.repo.Update(ctx, buyer); err != nil {
		return nil, err
	}

	return buyer, nil
}
