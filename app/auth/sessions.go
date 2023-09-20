package auth

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
