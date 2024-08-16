package user

import (
	"bytes"
)

type User struct {
	id           Id
	email        string
	name         string
	admin        bool
	passwordHash []byte
}

func New(email, name string) *User {
	return &User{
		id:    newId(),
		email: email,
		name:  name,
	}
}

func Restore(id Id, email, name string, admin bool, passwordHash []byte) *User {
	return &User{
		id:           id,
		email:        email,
		name:         name,
		admin:        admin,
		passwordHash: passwordHash,
	}
}

func (user *User) Id() Id {
	return user.id
}

func (user *User) Email() string {
	return user.email
}

func (user *User) Name() string {
	return user.name
}

func (user *User) IsAdmin() bool {
	return user.admin
}

func (user *User) PasswordHash() []byte {
	return user.passwordHash
}

func (user *User) PromoteAdmin() {
	user.admin = true
}

func (user *User) DemoteAdmin() {
	user.admin = false
}

func (user *User) ChangeEmail(email string) {
	user.email = email
}

func (user *User) ChangeName(name string) {
	user.name = name
}

func (user *User) ChangePasswordHash(passwordHash []byte) {
	user.passwordHash = passwordHash
}

func (user *User) PasswordSet() bool {
	return len(user.passwordHash) > 0
}

func (user *User) Equal(other *User) bool {
	return user.id == other.id &&
		user.email == other.email &&
		user.name == other.name &&
		user.admin == other.admin &&
		bytes.Equal(user.passwordHash, other.passwordHash)
}
