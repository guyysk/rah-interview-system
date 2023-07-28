package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	database "github.com/rbh-interview-system/src"
	"github.com/rbh-interview-system/src/model"
	"github.com/rbh-interview-system/src/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	client := database.ConnectDB()
	db := client.Database("topic")
	userRepository := repository.NewUserRepository(db)
	topicRepository := repository.NewTopicRepository(db)

	r := gin.New()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/topics", func(ctx *gin.Context) {
		limit, _ := strconv.Atoi(ctx.Query("limit"))
		page, _ := strconv.Atoi(ctx.Query("page"))
		status := ctx.Query("status")
		filter := bson.M{}

		if status != "" {
			split := strings.Split(status, ",")
			filter["status"] = bson.M{"$in": split}
		}

		topics, err := topicRepository.GetTopics(ctx, filter, limit, page)
		if err != nil {
			if errors.Is(err, repository.ErrTopicNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": topics})
	})

	r.PUT("/topics/:id", func(ctx *gin.Context) {
		topicId := ctx.Param("id")
		var body model.Topic
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		topic, err := topicRepository.UpdateTopicById(ctx, topicId, body)
		if err != nil {
			if errors.Is(err, repository.ErrTopicNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": topic})
	})

	r.GET("/users/:id", func(ctx *gin.Context) {
		userId := ctx.Param("id")
		user, err := userRepository.GetUserByUserID(ctx, userId)
		if err != nil {
			if errors.Is(err, repository.ErrUserNotFound) {
				ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
				return
			}
			ctx.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"data": user})
	})
	r.Run()
}
