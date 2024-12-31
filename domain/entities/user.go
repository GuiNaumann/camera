package entities

import "bear/utils"

const (
	UserTypeMaster UserType = 1

	UserTypeIsFlat3 UserType = 2

	UserTypeIsFlat2 UserType = 3

	UserTypeIsFlat1 UserType = 4
)

type UserType int64

// User basic user data in the application
type User struct {
	//ID Identifier attribute
	ID int64 `json:"id"`

	Name string `json:"name"`

	// SocialName - The preferred name that user choose to be called.
	SocialName *string `json:"socialName"`

	//Document CPF
	Document *string `json:"document,omitempty"`

	Address *Address `json:"address,omitempty"`

	Email string `json:"email"`

	BirthDate *utils.Date `json:"birthDate,omitempty"`

	//Image it can be a base64 or a URL stored on database
	Image *string `json:"image"`

	//ImageURL is the URL of where image was saved
	ImageURL string `json:"imageURL"`

	//ImageB64 is a representation of image as string (base 64)
	ImageB64 string `json:"imageB64"`

	Password string `json:"password,omitempty"`

	PasswordConfirmation string `json:"passwordConfirmation,omitempty"`

	//PasswordNotHashed used to send e-mail to user with their new password
	//ATTENTION!! NOT USE THIS TO LOG OR SAVE ON DATABASE
	PasswordNotHashed string `json:"omitempty"`

	Biography *string `json:"biography,omitempty"`

	AcademicEducation *string `json:"academicEducation,omitempty"`

	IsActive bool `json:"isActive"`

	////UserType user type represented by a number
	UserType UserType `json:"userType"`

	////UserTypeLabel user type represented by a label.
	UserTypeLabel string `json:"userTypeLabel,omitempty"`

	StatusCode StatusCode `json:"status_code"`

	//ModifiedAt is the date when the User was modified for the last time.
	ModifiedAt *utils.DateTime `json:"lastChange,omitempty"`

	//CurrentUser true if is current user.
	CurrentUser bool `json:"currentUser,omitempty"`

	IsForeigner bool `json:"isForeigner,omitempty"`

	//Dcoins dcoins wallet
	Decoins int64 `json:"dcoins"`

	//FirstLogin first login
	FirstLogin bool `json:"isFirstLogin"`

	//BadgeCode is the time card from users
	BadgeCode *int64 `json:"badgeCode"`

	//MothersName is the name of mother from user
	MothersName string `json:"mothersName"`

	//IsImported is true when the user was imported
	IsImported bool `json:"isImported"`

	//PasswordModifiedAt is the last time that tha password was modified
	PasswordModifiedAt *utils.DateTime `json:"passwordModifiedAt"`

	//TermAccept is true when user accepts terms and privacy policies
	TermAccept bool `json:"needAcceptNewTerm"`

	//AttendanceAt is the date when user is attended
	//AttendanceAt []AttendanceAts `json:"attendanceAts"`
}

type Address struct {
	Street *string `json:"street,omitempty"`

	AddressNumber *uint32 `json:"addressNumber,omitempty"`

	//District Neighborhood
	District *string `json:"district,omitempty"`

	//FederalUnity Referes to the states of Brazil
	FederalUnity int64 `json:"federalUnity,omitempty"`

	//ZipCode CEP
	ZipCode *string `json:"zipCode,omitempty"`

	//City
	City *string `json:"city,omitempty"`

	StateID *int `json:"stateID"`
}

// IsMaster check if this user has master privileges
func (u User) IsMaster() bool {
	return u.UserType.IsMaster()
}

// IsMaster check if this user type has master privileges
func (u UserType) IsMaster() bool {
	return u == UserTypeMaster
}

func (u *User) IsFlat3() bool {
	return u != nil && u.UserType.IsFlat3()
}

// IsAdmin check if this user type has administrative privileges
func (u UserType) IsFlat3() bool {
	return u == UserTypeIsFlat3
}

func (u *User) IsFlat2() bool {
	return u != nil && u.UserType.IsFlat2()
}

func (u UserType) IsFlat2() bool {
	return u == UserTypeIsFlat2
}

func (u *User) IsFlat1() bool {
	return u != nil && u.UserType.IsFlat1()
}

func (u UserType) IsFlat1() bool {
	return u == UserTypeIsFlat1
}

func (u User) Exists() bool {
	return u.StatusCode.Exists()
}

func (u StatusCode) Exists() bool {
	return u == StatusExist
}
