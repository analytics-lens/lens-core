package controller

import (
	model "core/internal/model"
	logger "core/pkg/log"

	"net/http"

	"github.com/gin-gonic/gin"
)

// GetMessages godoc
//
//	@Summary		Get chat messages
//	@Description	Returns messages for the current chat
//	@Tags			messages
//	@Produce		json
//	@Security		CookieAuth
//	@Success		200		{object}	[]model.Message
//	@Failure		400		{object}	map[string]any
//	@Failure		401		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/chat [get]
func (c *Controller) GetMessages(context *gin.Context) {

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

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": messages,
		},
	)
}
