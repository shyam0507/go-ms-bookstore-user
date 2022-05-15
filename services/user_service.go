package services

import (
	"go-ms-bookstore-user/domain/users"
	cryptoutils "go-ms-bookstore-user/utils/crypto_utils"
	datautils "go-ms-bookstore-user/utils/data_utils"
	"go-ms-bookstore-user/utils/errors"
)

const (
	UserStatus = "active"
)

var (
	UserService UserServiceInterface = &usersService{}
)

type usersService struct{}

type UserServiceInterface interface {
	GetUser(int64) (*users.User, *errors.RestErr)
	CreateUser(users.User) (*users.User, *errors.RestErr)
	UpdateUser(bool, users.User) (*users.User, *errors.RestErr)
	DeleteUser(users.User) *errors.RestErr
	SearchUser(string) (users.Users, *errors.RestErr)
}

func (s *usersService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	user := users.User{Id: userId}
	err := user.Get()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService) CreateUser(user users.User) (*users.User, *errors.RestErr) {

	if err := user.ValidateUser(); err != nil {
		return nil, err
	}

	user.Status = UserStatus
	user.DateCreated = datautils.GetNowString()
	user.Password = cryptoutils.GetMD5(user.Password)
	err := user.Save()

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *usersService) UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {

	if err := user.ValidateUser(); err != nil {
		return nil, err
	}

	currentUser, err := UserService.GetUser(user.Id)

	if err != nil {
		return nil, err
	}

	if isPartial {
		if user.FirstName != "" {
			currentUser.FirstName = user.FirstName
		}
		if user.LastName != "" {
			currentUser.LastName = user.LastName
		}
		if user.Email != "" {
			currentUser.Email = user.Email
		}
	} else {
		currentUser.FirstName = user.FirstName
		currentUser.LastName = user.LastName
		currentUser.Email = user.Email
	}

	updateErr := currentUser.Update()
	if updateErr != nil {
		return nil, updateErr
	}

	return currentUser, nil
}

func (s *usersService) DeleteUser(user users.User) *errors.RestErr {
	currentUser, err := UserService.GetUser(user.Id)

	if err != nil {
		return err
	}

	deleteErr := currentUser.Delete()
	if deleteErr != nil {
		return deleteErr
	}

	return nil
}

func (s *usersService) SearchUser(status string) (users.Users, *errors.RestErr) {
	user := users.User{}

	users, err := user.Search(status)

	if err != nil {
		return nil, err
	}

	return users, nil
}
