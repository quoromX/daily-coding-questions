package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/quoromx/irondoj/backend/internal/database"
)

type Service struct {
	store  *database.Store
	secret string
	pepper string
}

func New(store *database.Store, secret, pepper string) *Service {
	return &Service{store: store, secret: secret, pepper: pepper}
}

func (s *Service) Register(handle, displayName, email, password string) (database.User, string, error) {
	if len(strings.TrimSpace(password)) < 8 {
		return database.User{}, "", errors.New("password must be at least 8 characters")
	}
	if handle == "" {
		handle = strings.Split(email, "@")[0]
	}
	if displayName == "" {
		displayName = handle
	}
	hash := s.Hash(password)
	user, err := s.store.CreateUser(handle, displayName, email, hash)
	if err != nil {
		return database.User{}, "", err
	}
	return user, s.Token(user.ID), nil
}

func (s *Service) Login(email, password string) (database.User, string, error) {
	user, ok := s.store.UserByEmail(email)
	if !ok || user.PasswordHash != s.Hash(password) {
		return database.User{}, "", errors.New("invalid email or password")
	}
	return user, s.Token(user.ID), nil
}

func (s *Service) UserFromToken(token string) (database.User, bool) {
	userID, ok := s.Verify(token)
	if !ok {
		return database.User{}, false
	}
	return s.store.User(userID)
}

func (s *Service) Hash(password string) string {
	sum := sha256.Sum256([]byte(s.pepper + ":" + password))
	return hex.EncodeToString(sum[:])
}

func (s *Service) Token(userID string) string {
	expiry := time.Now().UTC().Add(24 * time.Hour).Unix()
	payload := userID + "." + base64.RawURLEncoding.EncodeToString([]byte(time.Unix(expiry, 0).Format(time.RFC3339)))
	signature := sign(payload, s.secret)
	return payload + "." + signature
}

func (s *Service) Verify(token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", false
	}
	payload := parts[0] + "." + parts[1]
	if !hmac.Equal([]byte(sign(payload, s.secret)), []byte(parts[2])) {
		return "", false
	}
	expiryBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}
	expiry, err := time.Parse(time.RFC3339, string(expiryBytes))
	if err != nil || time.Now().UTC().After(expiry) {
		return "", false
	}
	return parts[0], true
}

func sign(payload, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(payload))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func RandomSecret() string {
	var bytes [16]byte
	_, _ = rand.Read(bytes[:])
	return hex.EncodeToString(bytes[:])
}
