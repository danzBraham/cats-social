package securities

type PasswordHasher interface {
	HashPassword(password string) (hashedPassword string, err error)
	VerifyPassword(hashedPassword, password string) error
}
