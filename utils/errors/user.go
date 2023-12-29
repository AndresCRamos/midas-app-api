package errors

import "fmt"

const (
	USER_ALREADY_EXISTS = "A user with id %s already exists"
)

const (
	user_repository_error = "UserRepository: %s"
	user_service_error    = "UserService: %s"
)

// UserRepositoryError struct
type UserRepositoryError struct {
	Err error
}

func (ure UserRepositoryError) Error() string {
	return fmt.Sprintf(user_repository_error, ure.Err.Error())
}

func (ure *UserRepositoryError) Wrap(err error) {
	ure.Err = err
}

func (ure *UserRepositoryError) Unwrap() error {
	return ure.Err
}

// UserServiceError struct
type UserServiceError struct {
	Err error
}

func (use UserServiceError) Error() string {
	return fmt.Sprintf(user_service_error, use.Err.Error())
}

func (use *UserServiceError) Wrap(err error) {
	use.Err = err
}

func (use *UserServiceError) Unwrap() error {
	return use.Err
}
