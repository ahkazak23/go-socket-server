package controllers

import (
	"fmt"
	"go-socket-server/models"
	"go-socket-server/services"
)

type UserController struct {
	userService *services.UserService
}

// NewUserController creates a new instance of UserController
func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// --- User Management ---

// Register allows a new user to register with a username, password, role, and status
func (uc *UserController) Register(username, password, role, status string) string {
	err := uc.userService.RegisterUser(username, password, role, status)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Registration successful!"
}

// Login allows a user to log in with a username and password
func (uc *UserController) Login(username, password string) (string, bool, error) {
	user, err := uc.userService.LoginUser(username, password)
	if err != nil {
		return "", false, err
	}
	isAdmin := user.Role == "admin" && user.Status == "approved"
	return user.Username, isAdmin, nil
}

// ViewProfile allows a user to view their profile
func (uc *UserController) ViewProfile(username string) string {
	user, err := uc.userService.FindUserByUsername(username)
	if err != nil {
		return "Error: " + err.Error()
	}
	// Format profile information
	return fmt.Sprintf("\nName: %s\nSurname: %s\nFav Animal: %s\nFav Movie: %s\nYear of Birth: %s\nCity of Birth: %s\nFootball Team: %s",
		user.Name, user.Surname, user.FavAnimal, user.FavMovie, user.YearOfBirth, user.CityOfBirth, user.FootballTeam)
}

// UpdateProfile allows a user to update their profile information
func (uc *UserController) UpdateProfile(username, name, surname, favAnimal, favMovie, yearOfBirth, city, footballTeam string) string {
	err := uc.userService.UpdateUserProfile(username, name, surname, favAnimal, favMovie, yearOfBirth, city, footballTeam)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Profile updated successfully!" +
		"\nReturning to main menu."
}

// --- Blog Management ---

// PostBlog allows a user to create a blog post with a title and text
func (uc *UserController) PostBlog(username, title, text string) string {
	err := uc.userService.CreateBlog(username, title, text)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Blog posted successfully!"
}

// DeleteBlog allows a user to delete their own blog post
func (uc *UserController) DeleteBlog(username, blogID string) string {
	err := uc.userService.DeleteBlog(username, blogID)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Blog deleted successfully!"
}

// --- Admin Management ---

// ViewUsers allows an admin to view all registered users
func (uc *UserController) ViewUsers() string {
	users := uc.userService.GetAllUsers()
	response := "Users:\n"
	for _, user := range users {
		response += fmt.Sprintf("- %s (Role: %s, Status: %s)\n", user.Username, user.Role, user.Status)
	}
	return response
}

// DeleteUser allows an admin to delete a user by their username
func (uc *UserController) DeleteUser(username string) string {
	err := uc.userService.DeleteUser(username)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "User deleted successfully!"
}

// ViewPendingApprovals allows an admin to see pending admin applications
func (uc *UserController) ViewPendingApprovals() string {
	users := uc.userService.GetPendingAdminApprovals()
	response := "Pending Admin Approvals:\n"
	for _, user := range users {
		response += fmt.Sprintf("- %s\n", user.Username)
	}
	return response
}

// ApproveAdmin allows an admin to approve a user's admin request
func (uc *UserController) ApproveAdmin(username string) string {
	err := uc.userService.ApproveAdminRequest(username)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Admin request approved for user: " + username
}

// RejectAdmin allows an admin to reject a user's admin request
func (uc *UserController) RejectAdmin(username string) string {
	err := uc.userService.RejectAdminRequest(username)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Admin request rejected for user: " + username
}

// ApplyForAdmin allows a user to apply for admin status
func (uc *UserController) ApplyForAdmin(username string) string {
	err := uc.userService.ApplyForAdmin(username)
	if err != nil {
		return "Error: " + err.Error()
	}
	return "Admin application submitted successfully."
}

func (s *UserController) GetBlogsByUser(username string) []models.Blog {
	return s.userService.GetBlogsByUser(username)
}
