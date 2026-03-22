package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func PlantUsers(db *sqlx.DB, users []testUser) {
	for _, u := range users {
		id := uuid.NewString()
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
		if err != nil {
			log.Fatalf("failed to hash password for %s: %v", u.Email, err)
		}

		query := `INSERT INTO users (id, username, email, hashed_password)
				  VALUES ($1, $2, $3, $4)
				  ON CONFLICT (email) DO NOTHING`

		result, err := db.Exec(query, id, u.Username, u.Email, string(hashedPassword))
		if err != nil {
			log.Fatalf("failed to insert user %s: %v", u.Email, err)
		}

		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			fmt.Printf("⏭  User %s (%s) already exists, skipping\n", u.Username, u.Email)
		} else {
			fmt.Printf("Created user %s (%s)\n", u.Username, u.Email)
		}
	}
}
