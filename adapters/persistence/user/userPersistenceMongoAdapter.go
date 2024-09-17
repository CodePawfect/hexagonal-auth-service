// Package persistence provides functionality for user data persistence using MongoDB.
package persistence

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
	"user-auth-hexagonal-architecture/internal/domain"
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
		"role":      "USER",
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

// FindUser retrieves a user from the MongoDB database by their username.
//
// This method queries the MongoDB collection for a user document matching the given username.
// If found, it constructs and returns a domain.User struct with the user's information.
//
// Parameters:
//   - username: A string representing the username of the user to find.
//   - password: A string representing the password of the user (Note: This parameter is currently unused in the method body).
//
// Returns:
//   - domain.User: A User struct containing the user's information if found.
//   - error: An error if the user is not found or if there's a database error.
//     The error will be "user not found" if no matching user document is found,
//     or "failed to load user: [specific error]" for other database errors.
//
// Note:
//   - The password parameter is currently not used in the method body. Consider removing it if it's not needed.
//   - This method assumes that the "username" and "password" fields in the MongoDB document are strings.
//     It will panic if these fields are of a different type.
func (u *UserPersistenceMongoAdapter) FindUser(username string) (domain.User, error) {
	var result bson.M
	err := u.collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, fmt.Errorf("user not found")
		}
		return domain.User{}, fmt.Errorf("failed to load user: %w", err)
	}

	user := domain.User{
		Username: result["username"].(string),
		Password: result["password"].(string),
	}

	return user, nil
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
