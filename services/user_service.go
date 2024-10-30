package services

import (
	"errors"
	"fmt"
	"go-socket-server/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *models.InMemoryUserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo *models.InMemoryUserRepository) *UserService {
	return &UserService{repo: repo}
}

// --- User Management ---

// RegisterUser registers a new user with the specified username, password, role, and status
func (s *UserService) RegisterUser(username, password, role, status string) error {
	// Hash the password before storing it
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("Error hashing password: %s", err)
	}

	user := models.User{
		Username: username,
		Password: string(hashedPassword), // Store the hashed password
		Role:     role,
		Status:   status, // Can be "approved" or "pending"
	}
	return s.repo.CreateUser(user)
}

// LoginUser verifies the username and password for login
func (s *UserService) LoginUser(username, password string) (models.User, error) {
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}

	// Compare the hashed password with the password provided
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return models.User{}, fmt.Errorf("Invalid password")
	}

	return user, nil
}

// UpdateUserProfile updates the profile of the user with the given information
func (s *UserService) UpdateUserProfile(username, name, surname, favAnimal, favMovie, yearOfBirth, city, footballTeam string) error {
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return err
	}

	// Update profile information
	user.Name = name
	user.Surname = surname
	user.FavAnimal = favAnimal
	user.FavMovie = favMovie
	user.YearOfBirth = yearOfBirth
	user.CityOfBirth = city
	user.FootballTeam = footballTeam

	// Save the updated user back to the repository
	return s.repo.UpdateUser(user)
}

// --- Blog Management ---

// CreateBlog allows a user to create a new blog with the specified title and text
func (s *UserService) CreateBlog(username, title, text string) error {
	_, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return err
	}
	return s.repo.CreateBlog(username, title, text)
}

// DeleteBlog allows a user to delete their own blog post
func (s *UserService) DeleteBlog(username, blogID string) error {
	return s.repo.DeleteBlog(username, blogID)
}

// GetBlogsByUser fetches all blogs written by a specific user
func (s *UserService) GetBlogsByUser(username string) []models.Blog {
	return s.repo.GetBlogsByUser(username)
}

// --- Admin Management ---

// GetAllUsers returns all users in the system
func (s *UserService) GetAllUsers() []models.User {
	return s.repo.GetAllUsers()
}

// DeleteUser removes a user from the system
func (s *UserService) DeleteUser(username string) error {
	return s.repo.DeleteUser(username)
}

// GetPendingAdminApprovals fetches users who have applied for admin but are pending approval
func (s *UserService) GetPendingAdminApprovals() []models.User {
	return s.repo.GetPendingAdmins()
}

// ApproveAdminRequest approves a user's request to become an admin
func (s *UserService) ApproveAdminRequest(username string) error {
	return s.repo.ApproveAdmin(username)
}

// RejectAdminRequest rejects a user's request to become an admin
func (s *UserService) RejectAdminRequest(username string) error {
	return s.repo.RejectAdmin(username)
}

// ApplyForAdmin marks a user's status as "pending admin approval"
func (us *UserService) ApplyForAdmin(username string) error {
	user, exists := us.repo.Users[username]
	if !exists {
		return errors.New("user does not exist")
	}

	if user.Status == "pending" && user.Role == "admin" {
		return errors.New("admin application already pending")
	}

	// Mark the user as applying for admin status
	user.Status = "pending"
	user.Role = "admin"
	us.repo.Users[username] = user
	return nil
}
func (s *UserService) FindUserByUsername(username string) (models.User, error) {
	user, err := s.repo.FindUserByUsername(username)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}
