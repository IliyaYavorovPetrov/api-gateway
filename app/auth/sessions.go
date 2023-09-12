package auth

// Role 0.3.0
type Role string

const (
	// Unknown 0.3.0
	Unknown Role = "Unknown"

	// User 0.3.0
	User Role = "User"

	// Admin 0.3.0
	Admin Role = "Admin"
)

// NewRoleFromString 0.3.0
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

// IsValid 0.3.0
func (r Role) IsValid() bool {
	switch r {
	case User, Admin:
		return true
	}

	return false
}

// Session 0.3.0
type Session struct {
	UserID        string `redis:"userID"`
	Username      string `redis:"username"`
	UserRole      string `redis:"userRole"`
	IsBlacklisted bool   `redis:"isBlacklisted"`
}

// NewSession 0.3.0
func NewSession(userID string, username string, userRole string, isBlacklisted bool) (Session, error) {
	role, err := NewRoleFromString(userRole)
	if err != nil {
		return Session{}, ErrNotValidUserRole
	}

	return Session{userID, username, string(role), isBlacklisted}, nil
}
