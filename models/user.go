package models

import (
	"github.com/go-xorm/xorm"
	"github.com/hani17/chtq/models/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User Model
type User struct {
	ID       int64
	UserName string
	Email    string
	Password string

	CreatedAt time.Time `xorm:"created"`
	UpdatedAt time.Time `xorm:"updated"`
	DeletedAt time.Time `xorm:"deleted"`

	IsActive bool
	IsAdmin  bool
}

// IsUserExist checks if Username is used
func IsUserExist(name string) (bool, error) {
	return x.Get(&User{UserName: name})
}

// IsEmailUsed checks if Email is used
func IsEmailUsed(email string) (bool, error) {
	return x.Get(&User{Email: email})
}

// CreateUser create new record
func CreateUser(u *User) (err error) {
	isExist, err := IsUserExist(u.UserName)
	if err != nil {
		return err
	} else if isExist {
		return ErrUserAlreadyExist{u.UserName}
	}

	isExist, err = IsEmailUsed(u.Email)
	if err != nil {
		return err
	} else if isExist {
		return ErrEmailAlreadyUsed{Email: u.Email}
	}

	u.Password, _ = u.EncodePassword()

	sess := x.NewSession()
	defer sess.Close()

	if err = sess.Begin(); err != nil {
		return err
	}

	if _, err = sess.Insert(u); err != nil {
		return err
	}

	return sess.Commit()
}

// EncodePassword returns Password Hash
func (u *User) EncodePassword() (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(pwd), nil
}

// GetUserByName returns a user by given name.
func GetUserByName(name string) (*User, error) {
	if len(name) == 0 {
		return nil, errors.UserNotExist{UserID: 0, Name: name}
	}
	u := &User{UserName: name}
	has, err := x.Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.UserNotExist{UserID: 0, Name: name}
	}
	return u, nil
}

// GetUserById returns a user by given ID.
func getUserByID(e *xorm.Engine, id int64) (*User, error) {
	u := new(User)
	has, err := e.Id(id).Get(u)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, errors.UserNotExist{UserID: id, Name: ""}
	}
	return u, nil
}

// GetUserByID returns the user object by given ID if exists.
func GetUserByID(id int64) (*User, error) {
	return getUserByID(x, id)
}

// UserSignIn validates user name and password
func UserSignIn(username, password string) (*User, error) {
	var user *User
	user = &User{UserName: username}

	hasUser, err := x.Get(user)
	if err != nil {
		return nil, err
	}

	if hasUser {
		if user.ValidatePassword(password) {
			return user, nil
		}
	}
	return nil, errors.UserNotExist{UserID: user.ID, Name: user.UserName}
}

// ValidatePassword checks if given password matches the one belongs to the user
func (u *User) ValidatePassword(password string) bool {
	matching := ComparePasswords([]byte(u.Password), []byte(password))
	return matching
}

// ComparePasswords compares hashed password with input password
func ComparePasswords(hashedPassword, password []byte) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	if err != nil {
		return false
	}
	return true
}
