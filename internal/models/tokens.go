package models

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"time"
)

type Token struct {
	PlainText string    `json:"token"`
	UserID    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

var ScopeAuthentication string = "authorization"

func GenerateToken(userId int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token, nil
}

func (m *DatabaseModel) InsertToken(t *Token, u *User) error {
	stmt := `DELETE FROM token WHERE user_id = $1;`

	_, err := m.DB.Exec(stmt, u.ID)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO 
    		tokens 
    		(user_id, name, email, 
    		 token_hash, expiry, created_at)
			VALUES (
			        $1, $2, $3, 
			        $4, $5, $6);`

	_, err = m.DB.Exec(
		stmt,
		u.Lastname,
		u.Email,
		t.Hash,
		t.Expiry,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByAccessToken
func (m *DatabaseModel) GetUserByAccessToken(token string) (*pb.User, error) {
	tokenHash := sha256.Sum256([]byte(token))
	var u *pb.User

	stmt := `
			select 
				u.id, u.firstname, u.lastname, u.email	
			from
			    users u 
			inner join tokens t on (u.id = t.userId)
			where 
			    t.tokenHash = $1
			and
			    t.expiry > $2;`

	err := m.DB.QueryRow(
		stmt,
		tokenHash,
		time.Now()).Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Email,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}
