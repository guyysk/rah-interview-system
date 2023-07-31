package repository

import (
	"context"
	"errors"

	"github.com/rbh-interview-system/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type userRepository struct {
	userCollection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) UserRepository {
	userCollection := db.Collection("users")
	return &userRepository{userCollection: userCollection}
}

func (r userRepository) GetUserByUserID(ctx context.Context, userId string) (model.User, error) {
	var out user

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := r.userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&out)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.User{}, ErrUserNotFound
		}
	}

	return out.toUserModel(), nil
}

type user struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Name  string             `bson:"name,omitempty"`
	Email string             `bson:"email,omitempty"`
}

func (in user) toUserModel() model.User {
	return model.User{
		ID:    in.ID.Hex(),
		Name:  in.Name,
		Email: in.Email,
	}
}
