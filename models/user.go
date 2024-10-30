package models

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// User struct represents a user with various profile attributes
type User struct {
	Username     string
	Password     string
	Role         string // "user" or "admin"
	Status       string // "approved" or "pending"
	Name         string
	Surname      string
	FavAnimal    string
	FavMovie     string
	YearOfBirth  string
	CityOfBirth  string
	FootballTeam string
}

// Blog struct represents a blog post with an associated author (user)
type Blog struct {
	ID     string // Unique ID for the blog
	Author string // Username of the blog's author
	Title  string
	Text   string
}

// InMemoryUserRepository represents the in-memory database for users and blogs with file persistence
type InMemoryUserRepository struct {
	Users map[string]User
	Blogs map[string]Blog
	file  string // file path to persist data
}

// NewInMemoryUserRepository initializes a new repository with in-memory maps for users and blogs, and loads data from a file
func NewInMemoryUserRepository(file string) *InMemoryUserRepository {
	repo := &InMemoryUserRepository{
		Users: make(map[string]User),
		Blogs: make(map[string]Blog),
		file:  file,
	}
	repo.loadFromFile()
	return repo
}

// loadFromFile loads the users and blogs from the specified JSON file
func (repo *InMemoryUserRepository) loadFromFile() {
	fileData, err := ioutil.ReadFile(repo.file)
	if err != nil {
		// If the file does not exist, it will be created later
		fmt.Println("No data file found, starting fresh.")
		return
	}
	err = json.Unmarshal(fileData, repo)
	if err != nil {
		fmt.Printf("Error reading data from file: %s\n", err)
	}
}

// saveToFile writes the current users and blogs data to the specified JSON file
func (repo *InMemoryUserRepository) saveToFile() {
	fileData, err := json.MarshalIndent(repo, "", "  ")
	if err != nil {
		fmt.Printf("Error saving data to file: %s\n", err)
		return
	}
	err = ioutil.WriteFile(repo.file, fileData, 0644)
	if err != nil {
		fmt.Printf("Error writing data to file: %s\n", err)
	}
}

// --- User Methods ---

// CreateUser adds a new user to the repository and saves the changes to the file
func (repo *InMemoryUserRepository) CreateUser(user User) error {
	if _, exists := repo.Users[user.Username]; exists {
		return fmt.Errorf("User already exists")
	}
	repo.Users[user.Username] = user
	repo.saveToFile() // Persist changes to the file
	return nil
}

// FindUserByUsername retrieves a user by their username
func (repo *InMemoryUserRepository) FindUserByUsername(username string) (User, error) {
	user, exists := repo.Users[username]
	if !exists {
		return User{}, fmt.Errorf("User not found")
	}
	return user, nil
}

// UpdateUser updates an existing user's profile and saves the changes to the file
func (repo *InMemoryUserRepository) UpdateUser(user User) error {
	repo.Users[user.Username] = user
	repo.saveToFile() // Persist changes to the file
	return nil
}

// DeleteUser removes a user from the repository and saves the changes to the file
func (repo *InMemoryUserRepository) DeleteUser(username string) error {
	if _, exists := repo.Users[username]; !exists {
		return fmt.Errorf("User not found")
	}
	delete(repo.Users, username)
	repo.saveToFile() // Persist changes to the file
	return nil
}

// GetAllUsers returns all users
func (repo *InMemoryUserRepository) GetAllUsers() []User {
	users := []User{}
	for _, user := range repo.Users {
		users = append(users, user)
	}
	return users
}

// GetPendingAdmins returns all users with pending admin status
func (repo *InMemoryUserRepository) GetPendingAdmins() []User {
	pending := []User{}
	for _, user := range repo.Users {
		if user.Status == "pending" && user.Role == "admin" {
			pending = append(pending, user)
		}
	}
	return pending
}

// ApproveAdmin approves a user's admin request and saves the changes to the file
func (repo *InMemoryUserRepository) ApproveAdmin(username string) error {
	user, exists := repo.Users[username]
	if !exists {
		return fmt.Errorf("User not found")
	}
	if user.Status != "pending" {
		return fmt.Errorf("User is not pending approval")
	}
	user.Status = "approved"
	repo.Users[username] = user
	repo.saveToFile() // Persist changes to the file
	return nil
}

// RejectAdmin rejects a user's admin request and removes the user from the repository, saving the changes to the file
func (repo *InMemoryUserRepository) RejectAdmin(username string) error {
	user, exists := repo.Users[username]
	if !exists {
		return fmt.Errorf("User not found")
	}
	if user.Status != "pending" {
		return fmt.Errorf("User is not pending approval")
	}
	delete(repo.Users, username)
	repo.saveToFile() // Persist changes to the file
	return nil
}

// --- Blog Methods ---

// CreateBlog adds a new blog to the repository and saves the changes to the file
func (repo *InMemoryUserRepository) CreateBlog(username, title, text string) error {
	blogID := generateBlogID() // A function to generate a unique ID for the blog
	blog := Blog{
		ID:     blogID,
		Author: username, // Link the blog to the user who wrote it
		Title:  title,
		Text:   text,
	}
	repo.Blogs[blogID] = blog
	repo.saveToFile() // Persist changes to the file
	return nil
}

// DeleteBlog allows a user to delete their own blog and saves the changes to the file
func (repo *InMemoryUserRepository) DeleteBlog(username, blogID string) error {
	blog, exists := repo.Blogs[blogID]
	if !exists {
		return fmt.Errorf("Blog not found")
	}
	if blog.Author != username {
		return fmt.Errorf("You are not the author of this blog")
	}
	delete(repo.Blogs, blogID)
	repo.saveToFile() // Persist changes to the file
	return nil
}

// GetBlogsByUser returns all blogs written by a specific user
func (repo *InMemoryUserRepository) GetBlogsByUser(username string) []Blog {
	userBlogs := []Blog{}
	for _, blog := range repo.Blogs {
		if blog.Author == username {
			userBlogs = append(userBlogs, blog)
		}
	}
	return userBlogs
}

// --- Helper Functions ---

// generateBlogID generates a unique ID for each blog
func generateBlogID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// GetUserBlogs retrieves all blogs created by the specified user.
func (repo *InMemoryUserRepository) GetUserBlogs(username string) []Blog {
	userBlogs := []Blog{}
	for _, blog := range repo.Blogs {
		if blog.Author == username {
			userBlogs = append(userBlogs, blog)
		}
	}
	return userBlogs
}
