package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	// The file where emails will be stored - can be configured via environment variable
	emailFilePath string
	// The secret token for the protected endpoint - loaded from environment variable
	secretToken string
	host        string
)

var (
	// emailSet acts as our in-memory set for O(1) lookups.
	// map[string]struct{} is the idiomatic way to create a set in Go.
	emailSet = make(map[string]struct{})
	// mutex protects concurrent access to both emailSet and emailFilePath
	mutex = &sync.Mutex{}
)

// EmailRequest is the expected JSON structure for the POST request.
type EmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func main() {
	// 0. Load configuration from environment variables
	secretToken = os.Getenv("SECRET_TOKEN")
	if secretToken == "" {
		log.Fatal("ERROR: SECRET_TOKEN environment variable is not set")
	}

	host = os.Getenv("HOST")
	if host == "" {
		log.Fatal("ERROR: SECRET_TOKEN environment variable is not set")
	}

	// Set email file path from environment variable, default to "emails.txt"
	emailFilePath = os.Getenv("EMAIL_FILE_PATH")
	if emailFilePath == "" {
		emailFilePath = "emails.txt"
	}
	log.Printf("Using email file path: %s", emailFilePath)

	// 1. Load existing emails from the file into the in-memory set on startup.
	if err := loadEmailsFromFile(); err != nil {
		log.Printf("Warning: Could not load emails from %s: %v. Starting with an empty set.", emailFilePath, err)
	}
	log.Printf("Service started, loaded %d emails.", len(emailSet))

	// 2. Set up the Gin router.
	router := gin.Default()
	setupCors(router, host)

	// 3. Define endpoints.
	// POST /coming-soon (public)
	router.POST("/coming-soon", postEmailHandler)

	// GET /coming-soon (protected by middleware)
	// We create a group to apply the middleware only to specific routes.
	protected := router.Group("/")
	protected.Use(authMiddleware())
	{
		protected.GET("/coming-soon", getEmailsHandler)
	}

	// 4. Run the server.
	log.Println("Server starting on http://localhost:8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

// --- Handlers ---

// postEmailHandler handles new email submissions.
func postEmailHandler(c *gin.Context) {
	var req EmailRequest

	// 1. Bind and validate the incoming JSON.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request. 'email' field is required."})
		return
	}

	// 2. Validate the email format using the standard library.
	if _, err := mail.ParseAddress(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email format."})
		return
	}

	// 3. Normalize the email to lowercase.
	email := strings.ToLower(req.Email)

	// 4. Lock the mutex to ensure thread-safety for checking the map
	//    and writing to the file.
	mutex.Lock()
	defer mutex.Unlock()

	// 5. Check if the email already exists in the set.
	if _, exists := emailSet[email]; exists {
		c.JSON(http.StatusConflict, gin.H{"message": "Email already registered."})
		return
	}

	// 6. Add the new email to the set.
	emailSet[email] = struct{}{}

	// 7. Update the file on disk with all records from the map.
	if err := saveEmailsToFile(); err != nil {
		// If saving fails, roll back the in-memory change for consistency.
		delete(emailSet, email)
		log.Printf("ERROR: Failed to save email file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save email. Please try again."})
		return
	}

	// 8. Success!
	c.JSON(http.StatusCreated, gin.H{"message": "Email registered successfully."})
}

// getEmailsHandler returns the list of all registered emails.
func getEmailsHandler(c *gin.Context) {
	// Lock for read to prevent concurrent map modification while iterating.
	mutex.Lock()
	defer mutex.Unlock()

	// Convert the set (map keys) into a list (slice).
	emails := make([]string, 0, len(emailSet))
	for email := range emailSet {
		emails = append(emails, email)
	}

	c.JSON(http.StatusOK, emails)
}

// --- Middleware ---

// authMiddleware checks for the presence and correctness of the secret header.
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("X-Secret-Token")

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing X-Secret-Token header"})
			return
		}

		if token != secretToken {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Invalid secret token"})
			return
		}

		// Token is valid, proceed to the handler.
		c.Next()
	}
}

// --- File I/O Helpers ---

// loadEmailsFromFile reads the email file from disk into the in-memory set.
func loadEmailsFromFile() error {
	// Lock to prevent concurrent access during initialization.
	mutex.Lock()
	defer mutex.Unlock()

	file, err := os.Open(emailFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, which is fine.
		}
		return err // Other error (e.g., permissions)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := strings.TrimSpace(scanner.Text())
		if email != "" {
			// Add the normalized email to the set
			emailSet[strings.ToLower(email)] = struct{}{}
		}
	}
	return scanner.Err()
}

// saveEmailsToFile rewrites the entire email file with the current set.
// It assumes the mutex is already held by the caller.
func saveEmailsToFile() error {
	// Open the file with options: Write-only, Create if not exist, Truncate (clear) on open.
	file, err := os.OpenFile(emailFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %w", err)
	}
	defer file.Close()

	// Use a buffered writer for efficiency.
	writer := bufio.NewWriter(file)

	for email := range emailSet {
		// Write each email followed by a newline.
		if _, err := writer.WriteString(email + "\n"); err != nil {
			return fmt.Errorf("could not write email to file: %w", err)
		}
	}

	// Flush the buffer to ensure all data is written to disk.
	return writer.Flush()
}

func setupCors(r *gin.Engine, host string) {
	r.Use(cors.New(cors.Config{
		// 🚨 Allow your frontend origin
		AllowOrigins: []string{host},

		// Specify which methods are allowed
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// Specify which headers are allowed
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},

		// If you want to expose some headers to the browser
		ExposeHeaders: []string{"Content-Length"},

		// Enable this if your frontend needs to send cookies
		AllowCredentials: true,

		// How long the result of a preflight request can be cached
		MaxAge: 12 * time.Hour,
	}))
}
