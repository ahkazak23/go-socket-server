package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-socket-server/controllers"
	"go-socket-server/models"
	"go-socket-server/services"
)

var (
	connectedUsers = make(map[string]net.Conn)
	mu             sync.Mutex
)

func main() {
	// Initialize the in-memory user repository and services
	userRepo := models.NewInMemoryUserRepository("users.json")
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Start the server
	startServer(userController)
}

func startServer(controller *controllers.UserController) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("Server is listening on port 8080...")

	// Routine to display the count of connected users every 10 seconds
	go func() {
		for {
			mu.Lock()
			fmt.Printf("Connected users (%d):\n", len(connectedUsers))
			for ip := range connectedUsers {
				fmt.Println(ip)
			}
			mu.Unlock()
			time.Sleep(10 * time.Second)
		}
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConnection(conn, controller)
	}
}

func handleConnection(conn net.Conn, controller *controllers.UserController) {
	// Add the new connection to the map
	mu.Lock()
	connectedUsers[conn.RemoteAddr().String()] = conn
	mu.Unlock()

	defer func() {
		// Remove the connection from the map when the client disconnects
		mu.Lock()
		delete(connectedUsers, conn.RemoteAddr().String())
		mu.Unlock()
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var loggedInUser string
	var isAdmin bool

	// Display welcome message and prompt for login or registration
	writer.WriteString("******Welcome to the Go Socket Server!******\n")
	writer.WriteString("Type 'reg <username> <password>' to register.\n")
	writer.WriteString("Type 'log <username> <password>' to log in.\n")
	writer.Flush()

	for {
		// Read client input
		command, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from connection: %s\n", err)
			return
		}
		command = strings.TrimSpace(command)

		commandParts := strings.Fields(command)
		if len(commandParts) < 1 {
			writer.WriteString("Invalid command format.\n")
			writer.Flush()
			continue
		}

		var response string
		cmd := commandParts[0]

		switch cmd {
		case "reg":
			// Registration is allowed without login
			if len(commandParts) != 3 {
				response = "Usage: reg <username> <password>\n"
			} else {
				username := commandParts[1]
				password := commandParts[2]
				response = controller.Register(username, password, "user", "pending")
			}

		case "log":
			// User login
			if len(commandParts) != 3 {
				response = "Usage: log <username> <password>\n"
			} else {
				username := commandParts[1]
				password := commandParts[2]
				loggedInUser, isAdmin, err = controller.Login(username, password)
				if err != nil {
					response = "Invalid username or password. Please try again.\n"
				} else {
					if isAdmin {
						writer.WriteString("Login successful. Welcome Admin, " + username + "!\n")
					} else {
						writer.WriteString("Login successful. Welcome, " + username + "!\n")
					}
					response = displayMenu(isAdmin)
					writer.WriteString(response + "\n")
					writer.Flush()

					// Allow user to perform other actions after login
					for {
						// Read the next command
						command, err := reader.ReadString('\n')
						if err != nil {
							log.Printf("Error reading from connection: %s\n", err)
							return
						}
						command = strings.TrimSpace(command)
						commandParts := strings.Fields(command)
						if len(commandParts) < 1 {
							writer.WriteString("Invalid command format.\nReturning to main menu.\n")
							writer.Flush()
							continue
						}
						cmd = commandParts[0]

						switch cmd {
						// --- Profile Management ---
						case "view-profile":
							response = controller.ViewProfile(loggedInUser)
							writer.WriteString(response + "\nWould you like to edit your profile? (yes/no): ")
							writer.Flush()

							editResponse, _ := reader.ReadString('\n')
							editResponse = strings.TrimSpace(editResponse)

							if editResponse == "yes" {
								// Prompt for profile update fields one by one
								writer.WriteString("Name: ")
								writer.Flush()
								name, _ := reader.ReadString('\n')
								name = strings.TrimSpace(name)

								writer.WriteString("Surname: ")
								writer.Flush()
								surname, _ := reader.ReadString('\n')
								surname = strings.TrimSpace(surname)

								writer.WriteString("Favorite Animal: ")
								writer.Flush()
								favAnimal, _ := reader.ReadString('\n')
								favAnimal = strings.TrimSpace(favAnimal)

								writer.WriteString("Favorite Movie: ")
								writer.Flush()
								favMovie, _ := reader.ReadString('\n')
								favMovie = strings.TrimSpace(favMovie)

								writer.WriteString("Year of Birth: ")
								writer.Flush()
								yearOfBirth, _ := reader.ReadString('\n')
								yearOfBirth = strings.TrimSpace(yearOfBirth)

								writer.WriteString("City of Birth: ")
								writer.Flush()
								city, _ := reader.ReadString('\n')
								city = strings.TrimSpace(city)

								writer.WriteString("Football Team: ")
								writer.Flush()
								footballTeam, _ := reader.ReadString('\n')
								footballTeam = strings.TrimSpace(footballTeam)

								// Pass the collected information to update the profile
								text := controller.UpdateProfile(loggedInUser, name, surname, favAnimal, favMovie, yearOfBirth, city, footballTeam)
								writer.WriteString(text + "\n")
								response = displayMenu(isAdmin)
								break
							} else {
								response = "Profile edit canceled.\nReturning to main menu.\n"
								displayMenu(isAdmin)
								break
							}

						// --- Blog Management ---
						case "my-blogs":
							response = "Your Blogs:\n"
							blogs := controller.GetBlogsByUser(loggedInUser)
							for i, blog := range blogs {
								response += fmt.Sprintf("%d. %s\n%s\n", i+1, blog.Title, blog.Text)
							}
							writer.WriteString(response + "\nWould you like to post a new blog or delete one? (post/delete/exit): ")
							writer.Flush()

							actionResponse, _ := reader.ReadString('\n')
							actionResponse = strings.TrimSpace(actionResponse)

							switch actionResponse {
							case "post":
								writer.WriteString("Blog Title: ")
								writer.Flush()
								title, _ := reader.ReadString('\n')
								title = strings.TrimSpace(title)

								writer.WriteString("Blog Text: ")
								writer.Flush()
								text, _ := reader.ReadString('\n')
								text = strings.TrimSpace(text)

								response = controller.PostBlog(loggedInUser, title, text)
								displayMenu(isAdmin)
							case "delete":
								if len(blogs) == 0 {
									response = "You have no blogs to delete.\nReturning to main menu.\n"
								} else {
									writer.WriteString("Enter the blog number to delete: ")
									writer.Flush()
									indexInput, _ := reader.ReadString('\n')
									indexInput = strings.TrimSpace(indexInput)

									index, err := strconv.Atoi(indexInput)
									if err != nil || index < 1 || index > len(blogs) {
										response = "Invalid blog number.\nReturning to main menu.\n"
									} else {
										blogID := blogs[index-1].ID
										response = controller.DeleteBlog(loggedInUser, blogID)
									}
								}
							case "exit":
								response = displayMenu(isAdmin)
								break
							default:
								response = "Invalid option. Returning to main menu."
							}

						// --- Admin Management ---
						case "apply-admin":
							response = controller.ApplyForAdmin(loggedInUser)

						case "list-pending":
							if !isAdmin {
								response = "You do not have permission to perform this action.\nReturning to main menu.\n"
							} else {
								response = controller.ViewPendingApprovals()
								writer.WriteString(response + "\nWould you like to approve or reject any application? (approve/reject/exit): ")
								writer.Flush()

								approvalResponse, _ := reader.ReadString('\n')
								approvalResponse = strings.TrimSpace(approvalResponse)

								switch approvalResponse {
								case "approve":
									writer.WriteString("Username to approve: ")
									writer.Flush()
									username, _ := reader.ReadString('\n')
									username = strings.TrimSpace(username)
									response = controller.ApproveAdmin(username)
								case "reject":
									writer.WriteString("Username to reject: ")
									writer.Flush()
									username, _ := reader.ReadString('\n')
									username = strings.TrimSpace(username)
									response = controller.RejectAdmin(username)
								case "exit":
									response = "Exiting pending approvals.\nReturning to main menu."
									break
								default:
									response = "Invalid option.\n"

								}
							}
						case "list-users":
							if !isAdmin {
								response = "You do not have permission to perform this action.\nReturning to main menu.\n"
							} else {
								response = controller.ViewUsers()
								displayMenu(isAdmin)
							}
						case "exit":
							response = "Goodbye!"
							writer.WriteString(response + "\n")
							writer.Flush()
							return

						default:
							response = "Unknown command."
						}
						writer.WriteString(response + "\n")
						writer.Flush()
					}
				}
			}

		case "exit":
			response = "Goodbye!"
			writer.WriteString(response + "\n")
			writer.Flush()
			return

		default:
			response = "Unknown command.\n"
		}

		writer.WriteString(response + "\n")
		writer.Flush()
	}
}
func displayMenu(isAdmin bool) string {
	if isAdmin {
		return "Available commands:\n" +
			"- list-pending\n" +
			"- list-users\n" +
			"- view-profile\n" +
			"- my-blogs\n" +
			"- apply-admin\n" +
			"- exit\n"
	} else {
		return "Available commands:\n" +
			"- view-profile\n" +
			"- my-blogs\n" +
			"- apply-admin\n" +
			"- exit\n"
	}
}
