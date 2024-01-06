package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"greenlight.dzhdmitry.net/internal/validator"
)

const (
	ScopeActivation     = "activation"
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserId    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func generateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)

	if err != nil {
		return nil, err
	}

	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

func ValidateTokenPlaintext(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "plaintext", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "plaintext", "must be 26 bytes long")
}

type TokenRepository struct {
	DB *sql.DB
}

func (r TokenRepository) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userId, ttl, scope)

	if err != nil {
		return nil, err
	}

	err = r.Insert(token)

	return token, err
}

func (r TokenRepository) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
		VALUES ($1, $2, $3, $4)`

	args := []interface{}{token.Hash, token.UserId, token.Expiry, token.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := r.DB.QueryContext(ctx, query, args...)

	return err
}

func (r TokenRepository) DeleteAllForUser(scope string, userId int64) error {
	query := `DELETE FROM tokens 
		WHERE scope = $1 AND user_id = $2`

	args := []interface{}{scope, userId}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	_, err := r.DB.QueryContext(ctx, query, args...)

	return err
}

type MockTokenRepository struct {
	//
}

func (r MockTokenRepository) New(userId int64, ttl time.Duration, scope string) (*Token, error) {
	return nil, nil
}

func (r MockTokenRepository) Insert(token *Token) error {
	return nil
}

func (r MockTokenRepository) DeleteAllForUser(scope string, userId int64) error {
	return nil
}
