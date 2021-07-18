package formaterror

import (
	"strings"
)

var errorMessages = make(map[string]string)

var err error

// FormatError ...
func FormatError(err string) map[string]string {
	if strings.Contains(err, "nickname") {
		errorMessages["Taken_nickname"] = "Nickname Already Taken"
	}
	if strings.Contains(err, "email") {
		errorMessages["Taken_email"] = "Email Already Taken"
	}
	if strings.Contains(err, "title") {
		errorMessages["Taken_title"] = "Title Already Taken"
	}

	if strings.Contains(err, "hashedPassword") {
		errorMessages["incorrect_password"] = "Incorrect Password"
	}
	if strings.Contains(err, "record not found") {
		errorMessages["no_record"] = "User Not Found"
	}

	if strings.Contains(err, "double like") {
		errorMessages["Double_like"] = "You cannot like this post twice"
	}

	if len(errorMessages) > 0 {
		return errorMessages
	}

	if len(errorMessages) == 0 {
		errorMessages["Incorrect_details"] = "Incorrect Details"
		return errorMessages
	}
	return nil
}
