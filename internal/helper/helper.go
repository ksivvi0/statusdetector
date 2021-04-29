package helper

func CheckNullString(in string) bool {
	return len(in) > 0
}

func IsError(err error, needPanic bool) bool {
	if err != nil {
		if needPanic {
			panic(err)
		}
		return true
	}
	return false
}
