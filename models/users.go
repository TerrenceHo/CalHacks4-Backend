package models

import (
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name          string
	UserType      string
	Email         string `gorm:"not null;unique_index"`
	Password      string `gorm:"-"`
	PasswordHash  string `gorm:"not null"`
	PasswordReset bool
}

type UserDB interface {
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)

	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error
}

type UserService interface {
	Authenticate(email, password string) (*User, error)
	UserDB
}

func NewUserService(db *gorm.DB, pepper string) UserService {
	ug := &userGorm{db}
	uv := newUserValidator(ug, pepper)
	return &userService{
		UserDB: uv,
		pepper: pepper,
	}
}

// Checks if userService struct implements UserService interface
var _ UserService = &userService{}

type userService struct {
	UserDB
	pepper string
}

func (us *userService) Authenticate(email, password string) (*User, error) {
	foundUser, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(password+us.pepper))
	if err != nil {
		switch err {
		case bcrypt.ErrMismatchedHashAndPassword:
			return nil, ErrPasswordIncorrect
		default:
			return nil, err
		}
	}

	return foundUser, nil
}

// call a func type
// These functions that are of this type will run validation checks to on code
// to make sure they all comply with safety
type userValFunc func(*User) error

// When called, it will run all validation functions passed in, and if any
// return error then stop and return that error
func runUserValFuncs(user *User, fns ...userValFunc) error {
	for _, fn := range fns {
		err := fn(user)
		if err != nil {
			return err
		}
	}
	return nil
}

// Ensures that userValidator implements UserDB interface
var _ UserDB = &userValidator{}

// Struct that implements UserDB, which will include methods that run checks.
// pepper must be passed in
type userValidator struct {
	UserDB
	emailRegex *regexp.Regexp
	pepper     string
}

// Constructor for userValidator
func newUserValidator(udb UserDB, pepper string) *userValidator {
	return &userValidator{
		UserDB:     udb,
		emailRegex: regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,16}$`),
		pepper:     pepper,
	}
}

func (uv *userValidator) ByEmail(email string) (*User, error) {
	user := User{
		Email: email,
	}
	if err := runUserValFuncs(&user, uv.normalizeEmail, uv.emailFormat); err != nil {
		return nil, err
	}
	return uv.UserDB.ByEmail(user.Email)
}

func (uv *userValidator) Create(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordRequired,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Create(user)
}

func (uv *userValidator) Update(user *User) error {
	err := runUserValFuncs(user,
		uv.passwordMinLength,
		uv.bcryptPassword,
		uv.passwordHashRequired,
		uv.normalizeEmail,
		uv.requireEmail,
		uv.emailFormat,
		uv.emailIsAvail)
	if err != nil {
		return err
	}
	return uv.UserDB.Update(user)
}

func (uv *userValidator) Delete(id uint) error {
	var user User
	user.ID = id
	err := runUserValFuncs(&user, uv.idGreaterThan(0))
	if err != nil {
		return err
	}
	return uv.UserDB.Delete(id)
}

// bcryptPassword sets the password to empty string and password hash to the
// hash(password + salt + pepper)
func (uv *userValidator) bcryptPassword(user *User) error {
	if user.Password == "" {
		return nil
	}
	pwBytes := []byte(user.Password + uv.pepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(pwBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return nil
}

// n is usually zero.  Prevents this from searching database from an id that is
// less than 1.
func (uv *userValidator) idGreaterThan(n uint) userValFunc {
	return userValFunc(func(user *User) error {
		if user.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

// Normalizes emails by removing spaces and making all lowercase
func (uv *userValidator) normalizeEmail(user *User) error {
	user.Email = strings.ToLower(user.Email)
	user.Email = strings.TrimSpace(user.Email)
	return nil
}

func (uv *userValidator) requireEmail(user *User) error {
	if user.Email == "" {
		return ErrEmailRequired
	}
	return nil
}

func (uv *userValidator) emailFormat(user *User) error {
	if user.Email == "" {
		return nil
	}
	if !uv.emailRegex.MatchString(user.Email) {
		return ErrEmailInvalid
	}
	return nil
}

func (uv *userValidator) emailIsAvail(user *User) error {
	existing, err := uv.ByEmail(user.Email)
	if err == ErrEmailNotFound {
		// Email has not yet been taken
		return nil
	}

	// Otherwise, we have found a user with this email
	if user.ID != existing.ID {
		return ErrEmailTaken
	}
	return nil
}

func (uv *userValidator) passwordMinLength(user *User) error {
	if user.Password == "" {
		return nil
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	return nil
}

func (uv *userValidator) passwordRequired(user *User) error {
	if user.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

func (uv *userValidator) passwordHashRequired(user *User) error {
	if user.PasswordHash == "" {
		return ErrPasswordRequired
	}
	return nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db *gorm.DB
}

// Look up a user by the ID provided
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	if err == ErrResourceNotFound {
		return nil, ErrIDInvalid
	}
	return &user, err
}

func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	if err == ErrResourceNotFound {
		return nil, ErrEmailNotFound
	}
	return &user, err
}

// Take in a pointer to a user and create it in database
func (ug *userGorm) Create(user *User) error {
	return ug.db.Create(user).Error
}

// Update will update the provided user with all of the data
// in the provided user object.
func (ug *userGorm) Update(user *User) error {
	return ug.db.Model(user).Updates(user).Error
}

// Delete will delete the user with the provided ID
func (ug *userGorm) Delete(userID uint) error {
	var user User
	user.ID = userID
	return ug.db.Delete(user).Error
}

// Close closes the connection to database
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// Helper Functions

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrResourceNotFound
	}
	return err
}
