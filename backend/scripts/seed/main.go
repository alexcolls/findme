package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load environment
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("🌱 Seeding database...")

	// Seed users
	if err := seedUsers(db); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ Database seeded successfully!")
}

func seedUsers(db *sql.DB) error {
	users := []struct {
		email       string
		password    string
		fullName    string
		dateOfBirth string
		gender      string
		bio         string
	}{
		{
			email:       "john@example.com",
			password:    "Password123!",
			fullName:    "John Doe",
			dateOfBirth: "1995-06-15",
			gender:      "male",
			bio:         "Adventure seeker and coffee enthusiast. Love hiking and exploring new places!",
		},
		{
			email:       "jane@example.com",
			password:    "Password123!",
			fullName:    "Jane Smith",
			dateOfBirth: "1993-03-20",
			gender:      "female",
			bio:         "Book lover and travel addict. Always up for spontaneous road trips!",
		},
		{
			email:       "alex@example.com",
			password:    "Password123!",
			fullName:    "Alex Johnson",
			dateOfBirth: "1997-09-10",
			gender:      "other",
			bio:         "Foodie and music lover. Let's discover the best restaurants in town together!",
		},
		{
			email:       "sarah@example.com",
			password:    "Password123!",
			fullName:    "Sarah Williams",
			dateOfBirth: "1996-01-25",
			gender:      "female",
			bio:         "Yoga instructor and wellness enthusiast. Looking for someone who values health!",
		},
		{
			email:       "mike@example.com",
			password:    "Password123!",
			fullName:    "Mike Brown",
			dateOfBirth: "1994-11-05",
			gender:      "male",
			bio:         "Tech geek and gamer. Netflix and chill? More like code and compile! 😄",
		},
	}

	for _, user := range users {
		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// Check if user exists
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", user.email).Scan(&exists)
		if err != nil {
			return err
		}

		if exists {
			log.Printf("  ⏭  User already exists: %s", user.email)
			continue
		}

		// Insert user
		_, err = db.Exec(`
			INSERT INTO users (email, password_hash, full_name, date_of_birth, gender, bio, verified, active)
			VALUES ($1, $2, $3, $4, $5, $6, true, true)
		`, user.email, hashedPassword, user.fullName, user.dateOfBirth, user.gender, user.bio)

		if err != nil {
			return fmt.Errorf("failed to create user %s: %v", user.email, err)
		}

		log.Printf("  ✓ Created user: %s (password: %s)", user.email, user.password)
	}

	return nil
}
