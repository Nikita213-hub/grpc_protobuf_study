package validators

// Its a temp approach i will add more comprehensive check + user storage later
func CheckPassword(password string) bool {
	secretPassword := "secret_password"
	return password == secretPassword
}
