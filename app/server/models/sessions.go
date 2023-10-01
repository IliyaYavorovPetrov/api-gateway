package models

import (
	"encoding/json"
	"fmt"
)

type Role string

const (
	Unknown Role = "Unknown"
	User    Role = "User"
	Admin   Role = "Admin"
)

func NewRoleFromString(s string) (Role, error) {
	switch s {
	case string(User):
		return User, nil
	case string(Admin):
		return Admin, nil
	default:
		return Unknown, ErrNotValidUserRole
	}
}

func (r Role) IsValid() bool {
	switch r {
	case User, Admin:
		return true
	}

	return false
}

type Session struct {
	UserID        string `redis:"userID" json:"userID"`
	Username      string `redis:"username" json:"username"`
	UserRole      string `redis:"userRole" json:"userRole"`
	IsBlacklisted bool   `redis:"isBlacklisted" json:"isBlacklisted"`
}

func NewSession(userID string, username string, userRole string, isBlacklisted bool) (Session, error) {
	role, err := NewRoleFromString(userRole)
	if err != nil {
		return Session{}, ErrNotValidUserRole
	}

	return Session{userID, username, string(role), isBlacklisted}, nil
}

func (s Session) Equals(other interface{}) bool {
	if session, ok := other.(Session); ok {
		return s == session
	}

	return false
}

func (s Session) ToString() string {
	data, err := json.Marshal(s)
	if err != nil {
		fmt.Println("error encoding to json:", err)
		return ""
	}

	return string(data)
}
