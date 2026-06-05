package auth

import (
	"fmt"
	//"hash"
	"net/mail"
	"time"

	"github.com/MdRasB/LogLine/internal/db"
	"github.com/MdRasB/LogLine/internal/model"
	"github.com/google/uuid"
)

type Service struct{
	users *db.UserStore
	sessions *db.SessionStore
}

func NewService(
	users *db.UserStore,
	sessions *db.SessionStore,
) *Service {
	return &Service{
		users:    users,
		sessions: sessions,
	}
}

func (s *Service) Register(req RegisterRequest) error {

	if err := ValidateEmail(req.Email); err != nil {
		return fmt.Errorf("email validation failed: %w", err)
	}

	hash, err := HashPassword(req.Password)
	if err != nil {
		return err 
	}

	user := model.User{
		ID: uuid.New(),
		Email: req.Email,
		PasswordHash: hash,
		CreatedAt: time.Now(),
	}

	err = s.users.CreateUser(user)

	if err != nil {
		return fmt.Errorf("error user registration: %w", err)
	}

	return nil
}

func (s *Service) Login(req LoginRequest) (string, error) {
	if err := ValidateEmail(req.Email); err != nil {
		return "", fmt.Errorf("invalid email format: %w", err)
	}

	user, err := s.users.GetUserByEmail(req.Email)
	if err != nil {
		return "", err
	}

	err = VerifyPassword(req.Password, user.PasswordHash)
	if err != nil {
		return "", err 
	}

	token, tokenhash, err := GenerateSessionToken()
	if err != nil {
		return "", err
	}

	session := model.Session{
		ID: uuid.New(),
		UserID: user.ID,
		TokenHash: tokenhash,
		ExpiresAt: time.Now().Add(SessionDuration),
		CreatedAt: time.Now(),
	}

	err = s.sessions.CreateSession(session)
	if err != nil {
		return "", err 
	}

	return token, nil
}

func (s *Service) Logout(token string) error {
	if token == "" {
		return db.ErrInvalidToken
	}

	hash := HashSessionToken(token)

	session, err := s.sessions.GetSessionByTokenHash(hash)
	if err != nil {
		return err 
	}

	err = s.sessions.DeleteSession(session.ID)
	if err != nil {
		return err 
	}	

	return nil
}

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	
	return err 
}