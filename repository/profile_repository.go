package repository

import (
	"context"

	"github.com/Abenuterefe/a2sv-project/domain/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type profileRepository struct {
	collection *mongo.Collection
}

func NewProfileRepository(db *mongo.Database) *profileRepository {
	return &profileRepository{
		collection: db.Collection("profiles"), // assumes collection name is "profiles"
	}
}

// UpdateProfile updates bio and optionally profile picture without wiping other fields
func (r *profileRepository) UpdateProfile(ctx context.Context, profile *entities.Profile) error {
	// Get existing profile first
	existing, _ := r.FindByUserID(ctx, profile.UserID)
	if existing != nil {
		// Merge values so nothing gets overwritten with empty
		if profile.Bio == "" {
			profile.Bio = existing.Bio
		}
		if profile.ProfilePicture == "" {
			profile.ProfilePicture = existing.ProfilePicture
		}
	}

	filter := bson.M{"userId": profile.UserID}
	update := bson.M{
		"$set": bson.M{
			"bio":            profile.Bio,
			"profilePicture": profile.ProfilePicture,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}

// FindByUserID gets profile by userId
func (r *profileRepository) FindByUserID(ctx context.Context, userID string) (*entities.Profile, error) {
	filter := bson.M{"userId": userID}

	var profile entities.Profile
	err := r.collection.FindOne(ctx, filter).Decode(&profile)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // no profile yet
		}
		return nil, err
	}

	return &profile, nil
}

// UpdateProfilePicture updates profile picture without removing other fields
func (r *profileRepository) UpdateProfilePicture(ctx context.Context, userID string, picturePath string) error {
	// Get existing profile first
	existing, _ := r.FindByUserID(ctx, userID)
	bio := ""
	if existing != nil {
		bio = existing.Bio
	}

	filter := bson.M{"userId": userID}
	update := bson.M{
		"$set": bson.M{
			"profilePicture": picturePath,
			"bio":            bio, // keep existing bio
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := r.collection.UpdateOne(ctx, filter, update, opts)
	return err
}
