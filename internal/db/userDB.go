package db

import (
	"Database_Project/internal/structs"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
)

type UserDB struct {
	Client *sql.DB
}

// UserExists checks if a user with given username and password exists in the database.
func (db *UserDB) UserExists(username string, password string) (bool, error) {
	// Query to check if the user exists
	queryStmt := `SELECT COUNT(1) FROM UserAccount WHERE Username=? AND Password=?`

	// Use QueryRow to return a row and scan the returned id into the User struct
	var exists bool
	err := db.Client.QueryRow(queryStmt, username, password).Scan(&exists)

	if err != nil {
		// If an error is returned from the query, return a more contextual error message and the error
		return false, fmt.Errorf("unable to check if user exists: %v", err)
	}

	return exists, nil
}

// CheckLogin checks if the given username and password match a record in the database.
func (db *UserDB) CheckLogin(username string, password string) (structs.ActiveUser, error) {
	var user structs.ActiveUser

	queryStmt := `SELECT ID, Password FROM UserAccount WHERE Username=?`

	err := db.Client.QueryRow(queryStmt, username).Scan(&user.ID, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			// Username was not found
			return user, fmt.Errorf("username not found")
		}
		// Another error occurred
		return user, fmt.Errorf("internal error when fetching row: %v", err)
	}

	// this will check if the password hashes match, returns an error if they don't
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		// Password does not match stored hash
		return user, fmt.Errorf("unauthorized")
	}

	// Both username and password match the record in the database
	return user, nil
}

// RegisterUser creates a new user in the database.
func (db *UserDB) RegisterUser(userID, username, hashedPassword, email, firstName, lastName string, phone string) error {
	query := `
        INSERT INTO UserAccount (ID, Username, Password, Email, FirstName, LastName, Phone)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	_, err := db.Client.Exec(query, userID, username, hashedPassword, email, firstName, lastName, phone)
	if err != nil {
		log.Printf("Error inserting user into database: %v", err)
		return err
	}
	log.Printf("User successfully inserted into database. UserID: %s", userID)
	return nil
}

// GetUser retrieves a user with given username from the database.
func (db *UserDB) GetUser(username string) (structs.ActiveUser, error) {
	var user structs.ActiveUser

	query := `SELECT UserAccount.ID, UserAccount.Username, UserAccount.Email, UserAccount.Password, UserAccount.FirstName, UserAccount.LastName, UserAccount.Phone,
       UserAddress.Street, PostalCode.PostalCode 
		FROM UserAccount 
		LEFT JOIN UserAddress ON UserAccount.ID = UserAddress.UserAccountID
		LEFT JOIN PostalCode ON UserAddress.PostalCode = PostalCode.PostalCode
			WHERE UserAccount.Username = ?
`

	err := db.Client.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Address,
		&user.PostCode,
	)

	if err != nil {
		log.Printf("Error fetching user info for %v: %v", username, err)
		return user, err
	}

	return user, nil
}

// GetUser retrieves a user with given username from the database.
func (db *UserDB) GetUserByID(userID string) (structs.ActiveUser, error) {
	var user structs.ActiveUser

	query := `SELECT UserAccount.ID, UserAccount.Username, UserAccount.Email, UserAccount.Password, UserAccount.FirstName, UserAccount.LastName, UserAccount.Phone,
       UserAddress.Street, PostalCode.PostalCode 
		FROM UserAccount 
		LEFT JOIN UserAddress ON UserAccount.ID = UserAddress.UserAccountID
		LEFT JOIN PostalCode ON UserAddress.PostalCode = PostalCode.PostalCode
			WHERE UserAccount.ID = ?
`

	err := db.Client.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.Address,
		&user.PostCode,
	)

	if err != nil {
		log.Printf("Error fetching user info for %v: %v", userID, err)
		return user, err
	}

	return user, nil
}

// UpdateUserProfile updates the user's profile in the database, including their address.
func (db *UserDB) UpdateUserProfile(user structs.ActiveUser) error {
	// Begin transaction
	tx, err := db.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			// Rollback the transaction if an error occurred
			tx.Rollback()
			log.Printf("Transaction rolled back: %v", err)
			return
		}
		// Commit the transaction if everything is successful
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}()

	// Update UserAccount table
	_, err = tx.Exec(`
        UPDATE UserAccount SET Email=?, FirstName=?, LastName=?, Phone=? WHERE ID=?
    `, user.Email, user.FirstName, user.LastName, user.Phone, user.ID)
	if err != nil {
		return err
	}

	// Check if an address already exists for the user
	var existingAddressID string
	err = tx.QueryRow("SELECT ID FROM UserAddress WHERE UserAccountID = ?", user.ID).Scan(&existingAddressID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// If the user provided a new address, insert it
	if user.Address.String != "" && user.PostCode.String != "" {
		if existingAddressID != "" {
			// Update existing address
			_, err = tx.Exec(`
				UPDATE UserAddress
				SET Street=?, PostalCode=?
				WHERE UserAccountID=?
			`, user.Address.String, user.PostCode.String, user.ID)
			if err != nil {
				return err
			}
		} else {
			// Insert new address
			_, err = tx.Exec(`
				INSERT INTO UserAddress (ID, Street, PostalCode, UserAccountID)
				VALUES (?, ?, ?, ?)
			`, uuid.New().String(), user.Address.String, user.PostCode.String, user.ID)
			if err != nil {
				return err
			}
		}
	} else {
		// If the user did not provide a new address, delete the existing one if it exists
		if existingAddressID != "" {
			_, err = tx.Exec("DELETE FROM UserAddress WHERE ID = ?", existingAddressID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteUser deletes a user from the database along with their associated address and cart.
func (db *UserDB) DeleteUser(userID string) error {
	tx, err := db.Client.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("Transaction rolled back: %v", err)
			return
		}
		err = tx.Commit()
		if err != nil {
			log.Printf("Error committing transaction: %v", err)
		}
	}()

	// Delete Address associated with the user
	_, err = tx.Exec("DELETE FROM UserAddress WHERE UserAccountID = ?", userID)
	if err != nil {
		return err
	}

	// Delete User
	_, err = tx.Exec("DELETE FROM UserAccount WHERE ID = ?", userID) // Filter by ID
	if err != nil {
		return err
	}

	return nil
}

func (db *UserDB) UpdatePassword(userID string, newPassword string) error {
	query := `UPDATE UserAccount SET Password = ? WHERE ID = ?`
	_, err := db.Client.Exec(query, newPassword, userID)
	if err != nil {
		log.Printf("Error updating password for user %s: %v", userID, err)
		return err
	}
	return nil
}
