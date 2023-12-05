package models

import (
	"database/sql"
	pb "github.com/Chained/auth-service/github.com/Chained/auth-service"
	"github.com/Chained/auth-service/internal/driver"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type DatabaseModel struct {
	DB *sql.DB
}

type User struct {
	ID        int       `json:"id"`
	Firstname string    `json:"firstname"`
	Lastname  string    `json:"lastname"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

const (
	ErrUserNotFound  = "user was not found"
	ErrUserNotExists = "user does not exists"
	ErrWrongPassword = "invalid credentials"
)

/*
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    firstname VARCHAR(255) NOT NULL,
    lastname VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

*/

func NewDBModel() (*DatabaseModel, error) {
	db, err := driver.OpenDB()
	if err != nil {
		return nil, err
	}
	return &DatabaseModel{
		DB: db,
	}, nil
}

// GetUserByEmail returns user assiciated with given email
func (m *DatabaseModel) GetUserByEmail(email string) (*pb.User, error) {
	stmt := `
		SELECT 
		    * 
		FROM 
		    users 
		WHERE 
		    email = $1;`

	res := m.DB.QueryRow(stmt, email)
	if res.Err() != nil {
		return nil, res.Err()
	}

	var u *pb.User

	err := res.Scan(
		&u.Id,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.Role,
	)

	if err != nil {
		return nil, err
	}

	return u, nil
}

// CreateUser inserts user into database
func (m *DatabaseModel) CreateUser(
	firstname,
	lastname,
	email,
	password string) error {

	stmt := `INSERT INTO users(
            firstname,
            lastname,
        	email,
			password)
			VALUES (
			    $1, $2, 
			    $3, $4);`

	hashedPw, err := hashPassword(password)
	if err != nil {
		return err
	}

	row := m.DB.QueryRow(stmt,
		firstname,
		lastname,
		email,
		hashedPw,
	)

	if row.Err() != nil {
		return err
	}

	return nil
}

func (m *DatabaseModel) GetUser(id int) (*pb.User, error) {
	stmt := `SELECT * FROM users WHERE id = $1;`

	res, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	var u *pb.User

	if res.Next() {
		err = res.Scan(
			&u.Id,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Password,
		)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}

func (m *DatabaseModel) UpdateUser(request *pb.UpdateUserRequest) error {
	stmt := `
		UPDATE 
		    users 
		SET 
		    firstname = $1, 
		    lastname = $2, 
		    email = $3 
		WHERE id = $4;`

	res, err := m.DB.Exec(
		stmt,
		request.Firstname,
		request.Lastname,
		request.Email,
		request.Id)

	if err != nil {
		return err
	}

	if _, err = res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func (m *DatabaseModel) DeleteUser(id int) error {
	stmt := `DELETE FROM users WHERE id = $1`

	res, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	if _, err = res.RowsAffected(); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func comparePasswords(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}
