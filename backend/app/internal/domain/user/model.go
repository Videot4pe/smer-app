package user

import "time"

type User struct {
	Id         uint16    `json:"id" sql:"id"`
	Username   string    `json:"username" sql:"username"`
	Name       string    `json:"name" sql:"name"`
	Surname    string    `json:"surname" sql:"surname"`
	Patronymic string    `json:"patronymic" sql:"patronymic"`
	Email      string    `json:"email" validate:"required" sql:"email"`
	Password   string    `json:"password" validate:"required" sql:"password"`
	IsActive   bool      `json:"isActive" validate:"required" sql:"is_active"`
	IsVerified bool      `json:"isVerified" sql:"is_verified"`
	CreatedAt  time.Time `json:"createdAt" sql:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" sql:"updated_at"`

	AvatarId *uint16 `json:"avatarId" sql:"avatar_id"`
	Avatar   *string `json:"avatar"`
}

type VerificationData struct {
	Email     string    `json:"email" validate:"required" sql:"email"`
	Code      string    `json:"code" validate:"required" sql:"code"`
	ExpiresAt time.Time `json:"expiresAt" sql:"expires_at"`
	Type      string    `json:"type" sql:"type"`
}
