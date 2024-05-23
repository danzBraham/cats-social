package securities

type PasswordHasher interface {
	HashPassword(password string) (hashedPassword string, err error)
}
