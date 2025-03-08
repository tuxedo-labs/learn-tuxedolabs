package service

import (
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/internal/repository"
)

func GetUserByID(userID uint) (*entity.Users, error) {
	return repository.GetUserByID(userID)
}

func UpdateUserProfile(userID uint, updateUser entity.Users) error {
	user, err := repository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if updateUser.Name != "" {
		user.Name = updateUser.Name
	}
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != nil {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}
	if updateUser.Avatar != "" {
		user.Avatar = updateUser.Avatar
	}

	return repository.SaveUser(user)
}
