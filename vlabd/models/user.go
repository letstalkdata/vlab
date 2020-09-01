package models

import (
	"context"
	"fmt"
	"strings"
	"time"
	"os"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

type User struct {
	ID			uuid.UUID `json:"id"`
	CreatedAt	time.Time `json:"_"`
	UpdatedAt   time.Time `json:"_"`
	Email       string    `json:"email"`
	PasswordHash   string  `json:"_"`
	Password       string  `json:"password"`
	PasswordConfirm  string  `json:"password_confirm"`
}

func (u *User) Register(conn *pgx.Conn) error {
	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("password must be at least 4 characters long")
	}
	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("password do not match")
	}
	if len(u.Email) < 4 {
		return fmt.Errorf("Email must be at least 4 characters long")
	}

	u.Email = strings.ToLower(u.Email)
	row := conn.QueryRow(context.Background(), "SELECT id from user_account WHERE email = $1", u.Email)
	userLookup := User{}
	err := row.Scan(&userLookup)
	
	if err != pgx.ErrNoRows {
		fmt.Println("found user")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("An user with same email id already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("There was an error creating the account")
	}
	
	u.PasswordHash = string(pwdHash)
	//fmt.Println(u.PasswordHash)
	now := time.Now()
	_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, email, password_hash) VALUES($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash)
	return err
}

func (u *User) GetAuthToken() (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

func (u *User) IsAuthenticated(conn *pgx.Conn) error {
	row := conn.QueryRow(context.Background(), "SELECT id, password_hash FROM user_account WHERE email = $1", u.Email)
	err := row.Scan(&u.ID, &u.PasswordHash)
	
	if err == pgx.ErrNoRows {
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))

	if err != nil {
		return fmt.Errorf("Invalid login credentials")
	}

	return nil
}