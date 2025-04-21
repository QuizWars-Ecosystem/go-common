package jwt

type AuthLevel string

const (
	Unknown AuthLevel = "unknown"
	User    AuthLevel = "user"
	Admin   AuthLevel = "admin"
	Super   AuthLevel = "super"
)
