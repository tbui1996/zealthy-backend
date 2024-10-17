package services

import (
	"context"
	"log"

	"github.com/tbui1996/zealthy-backend/internal/core/domain"
	"github.com/tbui1996/zealthy-backend/internal/core/ports"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) ports.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) error {
	log.Printf("CreateUser service called with user: %+v", user)

	err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Printf("Error in repository Create method: %v", err)
		return err
	}

	log.Println("User created successfully in service")
	return nil
}

func (s *userService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

func (s *userService) GetAllUsers(ctx context.Context) ([]*domain.User, error) {
	return s.userRepo.GetAll(ctx)
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) error {
	log.Printf("UpdateUser service called with user: %+v", user)

	if user.Address != nil {
		log.Printf("Address details: Street: %v, City: %v, State: %v, Zip: %v",
			user.Address.Street, user.Address.City, user.Address.State, user.Address.Zip)
	} else {
		log.Println("Address is nil")
	}

	err := s.userRepo.Update(ctx, user)
	if err != nil {
		log.Printf("Error in repository Update method: %v", err)
		return err
	}

	log.Println("User updated successfully in service")
	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.userRepo.Delete(ctx, id)
}
