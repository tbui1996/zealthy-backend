package handlers

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tbui1996/zealthy-backend/internal/core/domain"
	"github.com/tbui1996/zealthy-backend/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	log.Println("CreateUser handler called")

	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Log the raw request body
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("Raw request body: %s", string(body))
	// Restore the request body for later use
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Received create user request for email: %s", input.Email)
	log.Printf("Password length: %d", len(input.Password)) // Don't log the actual password

	user := &domain.User{
		Email:    input.Email,
		Password: input.Password,
	}

	if err := h.userService.CreateUser(c.Request.Context(), user); err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	log.Printf("User created successfully with email: %s and ID: %s", user.Email, user.ID)

	// Don't return the password in the response
	user.Password = ""
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	log.Printf("Received update request for user ID: %s", id)

	var input domain.User
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("what is input.id: %s", input.ID)

	// Ensure the ID in the path matches the ID in the body
	if input.ID == "" {
		input.ID = id
	}

	log.Printf("Input data: %+v", input)

	if input.ID != id {
		log.Printf("ID mismatch: path ID %s, body ID %s", id, input.ID)
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID in path does not match ID in body"})
		return
	}

	// Fetch the existing user to ensure it exists
	existingUser, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		log.Printf("Error fetching user: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Merge the existing user data with the input data
	mergedUser := mergeUserData(existingUser, &input)

	// Perform the update with the merged user object
	if err := h.userService.UpdateUser(c.Request.Context(), mergedUser); err != nil {
		log.Printf("Error updating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	log.Printf("User updated successfully with ID: %s", mergedUser.ID)

	// Don't return the password in the response
	mergedUser.Password = ""
	c.JSON(http.StatusOK, mergedUser)
}

func mergeUserData(existing *domain.User, input *domain.User) *domain.User {
	merged := *existing // Create a copy of the existing user

	// Update fields only if they are provided in the input
	if input.Email != "" {
		merged.Email = input.Email
	}
	if input.Password != "" {
		merged.Password = input.Password
	}
	if input.AboutMe != nil {
		merged.AboutMe = input.AboutMe
	}
	if input.Birthdate != nil {
		merged.Birthdate = input.Birthdate
	}
	if input.Address != nil {
		if merged.Address == nil {
			merged.Address = &domain.Address{}
		}
		if input.Address.Street != nil {
			merged.Address.Street = input.Address.Street
		}
		if input.Address.City != nil {
			merged.Address.City = input.Address.City
		}
		if input.Address.State != nil {
			merged.Address.State = input.Address.State
		}
		if input.Address.Zip != nil {
			merged.Address.Zip = input.Address.Zip
		}
	}

	return &merged
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	user, err := h.userService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Users not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}
