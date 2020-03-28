package models

// User model containing Username, Email and Password
type Session struct {
	Username     string
	SessionToken string
	ExpiryTime   int64
}
