package entities

import (
	"camera/utils"
	"time"
)

const MinEntropyBits = 60

// LoginCredentials model to user makes the login request
type LoginCredentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthToken struct {
	Token string `json:"token"`
}

// UserPasswordRecoveryRequest store all the required data to start the password reset process
type UserPasswordRecoveryRequest struct {
	ID         int64     `json:"-"`
	UserID     int64     `json:"-"`
	Email      string    `json:"email"`
	Login      string    `json:"login"`
	UUID       string    `json:"-"`
	IPAddress  string    `json:"-"`
	Expiration time.Time `json:"-"`
	Status     uint8     `json:"-"`
}

// UserPasswordRecoveryResponse store the data sent to the user after a password recovery request
type UserPasswordRecoveryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// UserPasswordChangeRequest store all of the required data to change the user password
type UserPasswordChangeRequest struct {
	Login                   string `json:"-"`
	Token                   string `json:"-"`
	CurrentPassword         string `json:"current"`
	NewPassword             string `json:"password"`
	NewPasswordConfirmation string `json:"confirmation"`
	TermOfUse               bool   `json:"termOfUse"`
}

// UserCredentials store the credentials used to log a customer in the system
type UserCredentials struct {
	Login    string `json:"login,omitempty"`
	Password string `json:"password,omitempty"`
} // @name Usuario

// TokenRequestPasswordResetEnvelope used to get decrypted token from user reset password
type TokenRequestPasswordResetEnvelope struct {
	RequestUUID string `json:"requestUUID"`
	UserID      int64  `json:"userID"`
}

type TokenEnvelope struct {
	UserID     int64  `json:"userId"`
	DeviceUUID string `json:"deviceUUID"`
}

type RecoveryPasswordData struct {
	//Document CPF
	Document string `json:"document"`

	//Name Full name of user
	Name string `json:"name"`

	BirthDate *utils.Date `json:"birthDate"`

	//TimeCard is the time card from users
	TimeCard int64 `json:"badgeCode"`

	//MothersName is the name of mother from user
	MothersName string `json:"mothersName"`

	UUID string `json:"-"`

	IPAddress string `json:"-"`

	Expiration time.Time `json:"-"`

	Status uint8 `json:"-"`
}

type ResetPassword struct {
	Token                   string `json:"token"`
	NewPassword             string `json:"password"`
	NewPasswordConfirmation string `json:"confirmation"`
}

type NewPassword struct {
	NewPassword             string `json:"password"`
	NewPasswordConfirmation string `json:"confirmation"`
}
