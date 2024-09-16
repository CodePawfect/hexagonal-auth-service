// Package persistence provides functionality for user data persistence using MongoDB.
package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// UserPersistenceMongoAdapter implements the persistence layer for user-related operations.
// It encapsulates the MongoDB client and collection for user data.
type UserPersistenceMongoAdapter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserPersistenceMongoAdapter creates and initializes a new UserPersistenceMongoAdapter.
//
// It establishes a connection to MongoDB using the provided connection string and database name.
// The adapter uses a "user" collection within the specified database for all operations.
//
// Parameters:
//   - connectionString: MongoDB connection URI
//   - database: Name of the database to use
//
// Returns:
//   - *UserPersistenceMongoAdapter: A pointer to the newly created adapter
//   - error: An error if the connection fails or cannot be verified
func NewUserPersistenceMongoAdapter(client *mongo.Client, database string) (*UserPersistenceMongoAdapter, error) {
	collection := client.Database(database).Collection("user")

	return &UserPersistenceMongoAdapter{client, collection}, nil
}

// SaveUser stores user credentials in the MongoDB database.
//
// It creates a new document in the "user" collection with the provided username,
// hashed password, and the current timestamp.
//
// Parameters:
//   - username: The username of the user to be saved
//   - hashedPassword: The pre-hashed password of the user
//
// Returns:
//   - error: An error if the save operation fails, nil otherwise
//
// The function logs the ID of the newly inserted document on success.
func (u *UserPersistenceMongoAdapter) SaveUser(username string, hashedPassword string) error {
	user := bson.M{
		"username":  username,
		"password":  hashedPassword,
		"createdAt": time.Now(),
	}

	res, err := u.collection.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}

	log.Printf("User saved successfully with ID: %v", res.InsertedID)
	return nil
}

// IsUsernameAvailable checks if a given username is available for registration.
//
// It queries the database for an existing user with the provided username.
//
// Parameters:
//   - username: The username to check for availability
//
// Returns:
//   - bool: true if the username is available, false if it's already taken
//   - error: An error if the database query fails, nil otherwise
//
// Note: This function returns false for both an existing username and a database error.
// Check the error value to distinguish between these cases.
func (u *UserPersistenceMongoAdapter) IsUsernameAvailable(username string) (bool, error) {
	filter := bson.M{"username": username}
	existingUser := u.collection.FindOne(context.Background(), filter)
	if existingUser.Err() == nil {
		return false, nil
	}

	if !errors.Is(existingUser.Err(), mongo.ErrNoDocuments) {
		return false, existingUser.Err()
	}

	return true, nil
}

// LoadUser authenticates a user and generates a JWT token upon successful authentication.
//
// It takes a username and password as input, verifies the credentials against the stored
// information in the database, and returns a signed JWT token if authentication is successful.
//
// Parameters:
//   - username: The username of the user trying to authenticate.
//   - password: The password provided by the user for authentication.
//
// Returns:
//   - string: A signed JWT token if authentication is successful.
//   - error: An error if authentication fails or if there's an issue during the process.
//
// Possible errors:
//   - "failed to load user": If the user cannot be found in the database.
//   - "invalid password format in database": If the stored password is not in the expected format.
//   - "invalid password or username": If the provided credentials do not match.
//   - "error comparing passwords": If there's an unexpected error during password comparison.
//   - "error while creating jwt": If there's an issue generating the JWT token.
//
// Note: This function uses a hard-coded JWT key for demonstration purposes.
// In a production environment, the key should be securely stored and accessed.
func (u *UserPersistenceMongoAdapter) LoadUser(username string, password string) (string, error) {
	var result bson.M
	err := u.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("failed to load user: %w", err)
	}

	storedHash, ok := result["password"].([]byte)
	if !ok {
		return "", fmt.Errorf("invalid password format in database")
	}

	err = bcrypt.CompareHashAndPassword(storedHash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return "", fmt.Errorf("invalid password or username")
		}
		return "", fmt.Errorf("error comparing passwords: %w", err)
	}

	//TODO: Set claims for the JWT token (username, role 'USER')
	var jwtKey = []byte("my_secret_key") /*This is only for demo purpose, in a real application load the key from somewhere, for example an environment variable */
	token := jwt.New(jwt.SigningMethodHS256)
	signedString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error while creating jwt: %w", err)
	}

	return signedString, nil
}

// Close terminates the connection to the MongoDB database.
//
// It should be called when the UserPersistenceMongoAdapter is no longer needed to ensure
// proper cleanup of resources.
//
// Parameters:
//   - ctx: A context.Context for managing the lifecycle of the disconnect operation.
//
// Returns:
//   - error: An error if the disconnect operation fails, or nil if successful.
func (u *UserPersistenceMongoAdapter) Close(ctx context.Context) error {
	return u.client.Disconnect(ctx)
}
