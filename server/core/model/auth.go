package model

// AuthAccount represents an authenticated account that can be either a user or a circle.
// This separates authentication (who you are) from authorization (what account you're acting on behalf of).
type AuthAccount struct {
	// The authenticated user ID
	AuthUserId int64
	// The circle ID the user is acting on behalf of (optional)
	CircleId int64 `aip_pattern:"key=circle"`
	// The user ID the user is acting on behalf of (optional)
	UserId int64 `aip_pattern:"key=user"`
}
