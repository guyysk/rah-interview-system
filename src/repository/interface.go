package repository

import (
	"context"

	"github.com/rbh-interview-system/src/model"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository interface {
	GetUserByUserID(ctx context.Context, userId string) (model.User, error)
}

type TopicRepository interface {
	GetTopics(ctx context.Context, filter bson.M, limit int, page int) ([]model.Topic, error)
	UpdateTopicById(ctx context.Context, topicID string, updateTopic model.Topic) (model.Topic, error)
}
