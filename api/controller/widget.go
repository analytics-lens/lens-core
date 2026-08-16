package controller

import (
	"net/http"

	services "core/internal/services"
	logger "core/pkg/log"
	openai "core/pkg/openai"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WidgetRequest struct {
	Intent services.QueryIntent `json:"intent"`
}

// CreateWidget godoc
//
//	@Summary		Create dashboard widget
//	@Description	Builds a widget from a validated analytics intent
//	@Tags			widget
//	@Accept			json
//	@Produce		json
//	@Param			body	body	WidgetRequest	true	"Widget intent"
//	@Success		200		{object}	map[string]any
//	@Failure		400		{object}	map[string]any
//	@Failure		500		{object}	map[string]any
//	@Router			/v1/widget [post]
func (c *Controller) CreateWidget(context *gin.Context) {
	var request WidgetRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	query, err := services.BuildSQL(
		request.Intent,
	)

	if err != nil {
		context.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": err.Error(),
			},
		)

		return
	}

	var rows []map[string]any
	var queryResult *gorm.DB

	if len(query.Params) == 0 {
		queryResult = c.db.Raw(
			query.SQL,
		)
	} else {
		queryResult = c.db.Raw(
			query.SQL,
			query.Params,
		)
	}

	if 	err := queryResult.
		Scan(&rows).Error; 
		err != nil {
			
		logger.Error(
			"Failed to execute widget query: %v",
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

	fields := services.InferFields(rows)
	widget, err := openai.Client.GenerateWidget(
		request.Intent,
		fields,
		len(rows),
	)

	if err != nil {
		logger.Error(
			"Failed to generate widget: %v",
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

	result := map[string]any{
		"id":          uuid.New(),
		"type":        widget.Type,
		"title":       widget.Title,
		"description": widget.Description,
		"xKey":        widget.XKey,
		"dataKey":     widget.DataKey,
		"valueLabel":  widget.ValueLabel,
		"data":        rows,
	}

	logger.Info(
		"WIDGET: %+v",
		widget,
	)

	logger.Info(
		"ROWS: %+v",
		rows,
	)

	logger.Info(
		"FIELDS: %+v",
		fields,
	)

	context.JSON(
		http.StatusOK,
		gin.H{
			"data": result,
		},
	)
}
