package persistence

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

// UserPersistenceAdapter implements UserPersistencePort interface.
type UserPersistenceAdapter struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewUserPersistenceAdapter creates a new UserPersistenceAdapter struct and returns a pointer to it.
func NewUserPersistenceAdapter(connectionString string, database string) (*UserPersistenceAdapter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	collection := client.Database(database).Collection("user")

	return &UserPersistenceAdapter{client, collection}, nil
}

func (u *UserPersistenceAdapter) SaveUser(username string, hashedPassword string) error {
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

func (u *UserPersistenceAdapter) IsUsernameAvailable(username string) (bool, error) {
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

func (u *UserPersistenceAdapter) LoadUser(username string, password string) (string, error) {
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

	var jwtKey = []byte("my_secret_key") /*This is only for demo purpose, in a real application load the key from somewhere, for example an environment variable */
	token := jwt.New(jwt.SigningMethodHS256)
	signedString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", fmt.Errorf("error while creating jwt: %w", err)
	}

	return signedString, nil
}

func (u *UserPersistenceAdapter) Close(ctx context.Context) error {
	return u.client.Disconnect(ctx)
}
