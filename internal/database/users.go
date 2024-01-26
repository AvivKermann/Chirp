package database

import (
	"errors"
	"strings"

	"github.com/AvivKermann/Chirpy/models"
	"golang.org/x/crypto/bcrypt"
)

func (db *DB) CreateUser(email, password string) (models.ResponseUser, error) {

	dbContent, err := db.loadDB()
	if err != nil {
		return models.ResponseUser{}, err
	}

	for _, user := range dbContent.Users {
		if user.Email == email {
			return models.ResponseUser{}, errors.New("user already exist")
		}
	}

	index := len(dbContent.Users) + 1
	hashedPassword, err := HashPassword(password)

	if err != nil {
		return models.ResponseUser{}, nil
	}

	newUser := models.User{
		ResponseUser: models.ResponseUser{
			Email: email,
			ID:    index,
		},
		Password: hashedPassword,
	}
	dbContent.Users[index] = newUser

	err = db.writeDB(dbContent)
	if err != nil {
		return models.ResponseUser{}, err
	}

	return newUser.ResponseUser, nil
}

func (db *DB) UserLogin(email, password string) (models.ResponseUser, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return models.ResponseUser{}, err
	}

	isValid := ValidatePassword(user.Password, password)
	if !isValid {
		return models.ResponseUser{}, errors.New("password incorrect")
	}

	return user.ResponseUser, nil

}

func (db *DB) UpdateUser(email, password string, user models.User) (models.ResponseUser, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return models.ResponseUser{}, err
	}
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return models.ResponseUser{}, err
	}

	user.Email = email
	user.Password = hashedPassword

	dbContent.Users[user.ID] = user

	go db.writeDB(dbContent)

	return user.ResponseUser, nil
}
func (db *DB) GetUserByEmail(email string) (models.User, error) {
	dbContent, err := db.loadDB()
	if err != nil {
		return models.User{}, err
	}
	for _, user := range dbContent.Users {
		if user.Email == email {
			return user, nil
		}
	}
	return models.User{}, errors.New("email not found")

}

func (db *DB) GetUserById(userId int) (models.User, error) {
	dbContent, err := db.loadDB()

	if err != nil {
		return models.User{}, err
	}
	user, ok := dbContent.Users[userId]
	if !ok {
		return models.User{}, errors.New("user not found")
	}

	return user, nil
}

func ValidatePassword(hashedPassword []byte, password string) bool {
	err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	return err == nil
}

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return []byte{}, err
	}

	return hashedPassword, nil
}
func StripPrefix(header string) string {
	token := strings.TrimPrefix(header, "Bearer ")
	return token
}
