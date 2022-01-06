package user

type UserFormatter struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
	Token  string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:     user.ID,
		Name:   user.Name,
		Email:  user.Email,
		Avatar: user.AvatarFileName,
		Token:  token,
	}

	return formatter
}
