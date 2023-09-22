package models

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Loginresponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type UserProfile struct {
	FirstName  string `gorm:"not null" validate:"required,max=20"`
	MiddleName string
	LastName   string `gorm:"not null" validate:"required,max=20"`
	Mobile     string `gorm:"unique;not null" validate:"required,min=5,max=20"`
	Email      string `gorm:"unique;not null" validate:"required,min=5,max=45"`
	Intro      string
}
