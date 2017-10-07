package models

import "strings"

const (
	// ErrEmailNotFound is returned when an email cannot be found
	// in the database.
	ErrEmailNotFound modelError = "models: email not found"
	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"
	// ErrEmailRequired is returned when an email address is
	// not provided when creating a user
	ErrEmailRequired modelError = "models: email address is required"
	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid modelError = "models: email address is not valid"
	// ErrEmailTaken is returned when an update or create is attempted
	// with an email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"
	// ErrPasswordRequired is returned when a create is attempted
	// without a user password provided.
	ErrPasswordRequired modelError = "models: password is required"
	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user password that is less than 8 characters.
	ErrPasswordTooShort modelError = "models: password must be at least 8 characters long"
	// ErrVehicleRegNumNotFound is returned when looking for a vehicle
	// registration number that does not exist
	ErrVehicleRegNumNotFound modelError = `models: vehicle registration number not found.
		Please make sure you are the right user or have the correct vehicle registration number`
	// ErrVehicleIDInvalid is returned when querying for a vehicle id that
	// doesn't exist
	ErrVehicleIDInvalid modelError = "models: vehicle ID does not exist"
	// ErrVehicleRegNumRequired is returned when a vehicle does not have a
	// registration string
	ErrVehicleRegNumRequired modelError = "models: vehicle must have vegistration num"
	// ErrEmptyWeightRequired is returned when a vehicle entry is made without
	// an empty weight
	ErrEmptyWeightRequired modelError = "models: empty vehicle weight must not be empty"
	// ErrFullWeightRequired is returned when a vehicle entry is made without an
	// empty weight
	ErrFullWeightRequired modelError = "models: full vehicle weight must not be empty"
	// ErrPrevNotFilled is returned when trying to Post with same vehicleRegNum
	// twice
	ErrPrevNotFilled modelError = "models: previous vehicle with same registration number not completely filled in"
	// ErrPrevAlreadyFilled is returned when trying to PUT with the same vehicle
	// reg_num that was already updated once.
	ErrPrevAlreadyFilled modelError = "models: vehicle you are trying to update was already updated once.  You cannot update it twice."

	// privateError only for internal use only, not prod
	// ErrResourceNotFound is returned when a resource cannot be found in
	// database
	ErrResourceNotFound privateError = "models: resource not found"
	// ErrIDInvalid is returned when an invalid ID is provided
	// to a method like Delete.
	ErrIDInvalid privateError = "models: ID provided was invalid"
	// UserID needed
	ErrUserIDRequired privateError = "models: user ID is required"
	// ErrColumnNotFound is returned when looking up a column that doesn't exist
	ErrColumnNotFound privateError = "models: column doesn't exist"
)

// model error implements the error interface type, because it has the Error()
// method built in
type modelError string

func (e modelError) Error() string {
	return string(e)
}

func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

// privateError implements the error interface type, because it has the Error()
// method build in
type privateError string

func (e privateError) Error() string {
	return string(e)
}
