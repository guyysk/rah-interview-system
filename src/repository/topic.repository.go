package repository

import (
	"context"
	"errors"
	"time"

	"github.com/rbh-interview-system/src/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ErrTopicNotFound = errors.New("topic not found")
)

type topicRepository struct {
	topicCollection *mongo.Collection
}

func NewTopicRepository(db *mongo.Database) TopicRepository {
	topicCollection := db.Collection("topics")
	return &topicRepository{topicCollection: topicCollection}
}

func (r topicRepository) GetTopics(ctx context.Context, filter bson.M, limit int, page int) ([]model.Topic, error) {
	var topics []model.Topic
	var topic topic

	options := options.Find()
	options.SetLimit(int64(limit))
	options.SetSkip(int64(page) * int64(limit))

	cursor, err := r.topicCollection.Find(ctx, filter, options)
	if err != nil {
		defer cursor.Close(ctx)
		return nil, err
	}

	for cursor.Next(ctx) {
		err := cursor.Decode(&topic)
		if err != nil {
			return topics, err
		}
		topics = append(topics, topic.toTopicModel())
	}

	return topics, nil
}

func (r topicRepository) UpdateTopicById(ctx context.Context, topicID string, updateTopic model.Topic) (model.Topic, error) {

	var topic topic
	in := bson.M{}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	if updateTopic.Title != "" {
		in["title"] = updateTopic.Title
	}
	if updateTopic.Description != "" {
		in["description"] = updateTopic.Description
	}
	if updateTopic.Status != "" {
		in["status"] = updateTopic.Status
	}

	objId, _ := primitive.ObjectIDFromHex(topicID)

	err := r.topicCollection.FindOneAndUpdate(ctx, bson.M{"_id": objId}, bson.M{"$set": in}, &opt).Decode(&topic)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Topic{}, ErrTopicNotFound
		}
	}

	return topic.toTopicModel(), nil
}

type comment struct {
	Comment   string    `bson:"comment"`
	CreatedBy string    `bson:"createdBy"`
	CreatedAt time.Time `bson:"createdAt"`
}

type topic struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Title       string             `bson:"title,omitempty"`
	Description string             `bson:"description"`
	Comments    []comment          `bson:"comments"`
	Status      string             `bson:"status,omitempty"`
	CreatedBy   string             `bson:"createdBy"`
	CreatedAt   time.Time          `bson:"createdAt"`
}

func (in topic) toTopicModel() model.Topic {
	var comments []model.Comment
	for _, comment := range in.Comments {
		comments = append(comments, model.Comment{
			Comment:   comment.Comment,
			CreatedBy: comment.CreatedBy,
			CreatedAt: comment.CreatedAt,
		})
	}
	return model.Topic{
		ID:          in.ID.Hex(),
		Title:       in.Title,
		Description: in.Description,
		Comments:    comments,
		Status:      in.Status,
		CreatedBy:   in.CreatedBy,
		CreatedAt:   in.CreatedAt,
	}
}
