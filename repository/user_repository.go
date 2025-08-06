package repository

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		collection: db.Collection("user"),
		db:         db,
	}
}

// Create user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

// find user by their id
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// store token of user for authentication
func (r *userRepository) StoreToken(ctx context.Context, token *entities.Token) error {
	_, err := r.db.Collection("tokens").InsertOne(ctx, token)
	return err
}

// Find or check token if it exists in database
func (r *userRepository) FindToken(ctx context.Context, refreshToken string) (*entities.Token, error) {
	var token entities.Token
	err := r.db.Collection("tokens").FindOne(ctx, bson.M{"refresh_token": refreshToken}).Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// Delete token (expired token) from database
func (r *userRepository) DeleteToken(ctx context.Context, refreshToken string) error {
	_, err := r.db.Collection("tokens").DeleteOne(ctx, bson.M{"refresh_token": refreshToken})
	return err
}

// Find user by email verification token
func (r *userRepository) FindByVerificationToken(ctx context.Context, token string) (*entities.User, error) {
	var user entities.User
	err := r.collection.FindOne(ctx, bson.M{"verification_token": token}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update users verified status and verification token
func (r *userRepository) Update(ctx context.Context,user *entities.User) error{
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"verified": user.Verified,
		"verification_token": user.VerificationToken,
	}}

	_,err := r.collection.UpdateOne(ctx,filter,update)

	return err
}
