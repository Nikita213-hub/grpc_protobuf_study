package validators

func IsUserIdValid(userId string) bool {
	allowedUserIds := []string{"1", "2", "3"}
	isValid := false
	for _, v := range allowedUserIds {
		if userId == v {
			isValid = true
		}
	}
	return isValid
}
