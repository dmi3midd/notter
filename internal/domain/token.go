package domain

type Token struct {
	UserId       string `json:"userId" db:"user_id"`
	RefreshToken string `json:"refreshToken" db:"refresh_token"`
}

type TokenRepository interface {
	GetToken(refreshToken string) (*Token, error)
	Create(userId, refreshToken string) error
	Delete(refreshToken string) error
	Update(userId, refreshToken string) error
}
