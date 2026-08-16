package controller

import (
	"net/http"

	model "core/internal/model"
	services "core/internal/services"
	logger "core/pkg/log"
	openai "core/pkg/openai"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AskRequest struct {
	Query  string               `json:"query"`
	Intent services.QueryIntent `json:"intent"`
}

// Ask godoc
//
//	@Summary		Ask analytics agent
//	@Description	Progressively builds an analytics intent from the conversation
//	@Tags			agent
//	@Accept			json
//	@Produce		json
//	@Param			body	body	AskRequest	true	"Agent request"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/v1/ask [post]
func (c *Controller) Ask(context *gin.Context) {

	var request AskRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	if request.Query == "" {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "Missing query",
			},
		)
		return
	}

	var messages []model.Message
	if err := c.db.
		Order("created_at ASC").
		Find(&messages).Error; err != nil {

		logger.Error(
			"Failed to get messages: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	message := model.Message{
		ID: 		uuid.New(),
		Body:       request.Query,
		SenderRole: model.CHAT_ROLE_USER,
	}

	if err := c.db.Create(&message).Error; err != nil {

		logger.Error(
			"Failed to save user message: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	agent, err := openai.Client.GenerateIntent(
		request.Query,
		model.SupportSchema,
		request.Intent,
		messages,
	)

	if err != nil {
		logger.Error(
			"Failed to generate agent response: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	aiMessage := model.Message{
		ID: 		uuid.New(),
		Body:       agent.Reply,
		SenderRole: model.CHAT_ROLE_AI,
	}

	if err := c.db.Create(&aiMessage).Error; err != nil {

		logger.Error(
			"Failed to save AI message: %v",
			err,
		)

		context.JSON(
			http.StatusInternalServerError,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": gin.H{
				"status": agent.Status,
				"reply":  agent.Reply,
				"intent": agent.Intent,
			},
		},
	)
}
