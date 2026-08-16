package openai

import (
	"context"
	"core/internal/services"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const widgetPrompt = `
You are a dashboard widget planner.

Choose the best widget for
the user's request.

Available widgets:

- number
- bar
- line
- area
- progress

You are given ONLY the result schema.
You do not receive the underlying data.

Result fields:

%s

Return JSON only:

{
	"type": "number | bar | line | area | progress",
	"title": "string",
	"description": "string | null",
	"xKey": "string | null",
	"dataKey": "string | null",
	"valueLabel": "string | null"
	"max": "number | null"
}

Rules:

- Never invent field names.
- Use number for a single scalar result.
- Use bar for categorical comparisons.
- Use line for time-series trends.
- Use area for time-series magnitude/trend.
- Use progress only when the request
  clearly implies progress toward a target.
- xKey must be a field from the result.
- dataKey must be a numeric field.
- For number widgets, dataKey is REQUIRED.
- dataKey must be the numeric result field.
- For number widgets, xKey must be null.
- For bar, line, and area widgets, dataKey is REQUIRED.
- For bar, line, and area widgets, xKey is REQUIRED.
- Never return null for dataKey when a numeric field exists.
`

type WidgetResponse struct {
	Type        string   `json:"type"`
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	XKey        *string  `json:"xKey"`
	DataKey     *string  `json:"dataKey"`
	ValueLabel  *string  `json:"valueLabel"`
	Max         *float64 `json:"max"`
}

func (c *ClientStruct) GenerateWidget(
	userRequest services.QueryIntent,
	resultSchema any,
	rowCount int,
) (*WidgetResponse, error) {

	schemaJSON, err := json.MarshalIndent(
		resultSchema,
		"",
		"\t",
	)

	if err != nil {
		return nil, err
	}

	prompt := fmt.Sprintf(
		widgetPrompt,
		string(schemaJSON),
	)

	request := map[string]any{
		"request":      userRequest,
		"resultSchema": resultSchema,
		"rowCount":     rowCount,
	}

	requestJSON, err := json.Marshal(
		request,
	)

	if err != nil {
		return nil, err
	}

	response, err := c.Client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4_1Nano,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(
					prompt +
						"\n\nRequest:\n" +
						string(requestJSON),
				),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	var result WidgetResponse

	if 	err := json.Unmarshal(
			[]byte(response.OutputText()),
			&result,
		);
		err != nil {

		return nil, fmt.Errorf(
			"failed to parse widget response: %w\nresponse=%s",
			err,
			response.OutputText(),
		)
	}

	return &result, nil
}
