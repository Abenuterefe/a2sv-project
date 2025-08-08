package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
	resetTokenCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *userRepository {
	return &userRepository{
		db:         db,
		collection: db.Collection("user"),
		resetTokenCollection:db.Collection("reset_tokens"),
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
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": bson.M{
		"verified":           user.Verified,
		"verification_token": user.VerificationToken,
		"role":               user.Role,
	}}

	_, err := r.collection.UpdateOne(ctx, filter, update)

	return err
}

// Find user by their id repo func
func (r *userRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*entities.User, error) {
	var user entities.User

	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNilDocument {
			return nil, nil // if no user found
		}
		return nil, err //other db error
	}

	return &user, nil
}

// Deleting refresh token from cache or database
func (r *userRepository) DeleteTokenByUserID(ctx context.Context, userID string) error {
	filter := bson.M{"user_id": userID}
	_, err := r.db.Collection("tokens").DeleteOne(ctx, filter)
	return err
}


//RESET TOKEN REPOSITORY
//save reset token
func (r *userRepository) SaveResetToken (ctx context.Context, token *entities.ResetToken) error{
	_, err := r.resetTokenCollection.InsertOne(ctx, token)
	return err
}

//find resetToken by user id
func (r *userRepository) FindByResetToken(ctx context.Context, token string)(*entities.ResetToken,error){
	var resetToken entities.ResetToken
	err := r.resetTokenCollection.FindOne(ctx, map[string]interface{}{"token": token}).Decode(&resetToken)
	return &resetToken, err
}

//Delete reset token function
func (r *userRepository) DeleteResetToken(ctx context.Context, token string) error {
	_, err := r.resetTokenCollection.DeleteOne(ctx, map[string]interface{}{"token": token})
	return err
}

//Update user password
func (r *userRepository) UpdatePassword(ctx context.Context, userID primitive.ObjectID, hashedPassword string) error {
	filter := bson.M{"_id": userID}
	update := bson.M{
		"$set": bson.M{
			"password":  hashedPassword,
			"updatedAt": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}
