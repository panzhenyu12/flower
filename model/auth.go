package model

type LoginReq struct {
	UserName string `json:"userName" valid:"required"`
	Password string `json:"password" valid:"required"`
}

type UpdatePasswordRequest struct {
	OldPassword string
	NewPassword string
}

func EncryptPassword(password string) string {
	return password // _(:з」∠)_
}
