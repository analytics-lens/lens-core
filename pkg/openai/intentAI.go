package openai

import (
	"context"
	"core/internal/services"
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/responses"
)

const intentPrompt = `
You are an analytics conversation agent.

Your job is to progressively construct
a structured QueryIntent from the user's request.

You DO NOT write SQL.

You can only use metrics and dimensions
from this schema:

%s

Previous intent:

%s

Conversation:

%s

Return JSON only.

Response shape:

{
	"status": "clarifying | final",
	"reply": "string",
	"intent": {
		"metric": "metric name | null",
		"aggregation": "sum | avg | null",
		"dimension": "dimension name | null",
		"timeRange": {
			"type": "all | last_7_days | last_30_days",
			"field": "date | null"
		},
		"filters": [
			{
				"dimension": "dimension name",
				"operator": "eq",
				"value": "string"
			}
		]
	}
}

Rules:

- Never invent metrics.
- Never invent dimensions.
- Never invent filter fields.
- Only use aggregations allowed by the metric.
- Only use eq filters.
- If no dimension is required, use null.
- If the user asks for a trend over time,
  use date as the dimension.
- If the user asks for a comparison by channel,
  use channel as the dimension.
- If no time range is specified, use all.
- Preserve the previous intent when the user
  provides additional information.
- Do not discard information from the
  previous intent unless the user changes it.
- If the request is incomplete, return
  status "clarifying".
- Ask a concise clarification question
  in "reply".
- If the intent contains everything required
  to answer the request, return "final".
- Do not include SQL.
- Do not include database rows.
`

type IntentResponse struct {
	Status string               `json:"status"`
	Reply  string               `json:"reply"`
	Intent services.QueryIntent `json:"intent"`
}

func (c *ClientStruct) GenerateIntent(
	query string,
	schema any,
	previousIntent any,
	messages any,
) (*IntentResponse, error) {

	schemaJSON, err := json.MarshalIndent(
		schema,
		"", "\t",
	)

	if err != nil {
		return nil, err
	}

	previousIntentJSON, err := json.MarshalIndent(
		previousIntent,
		"", "\t",
	)

	if err != nil {
		return nil, err
	}

	messagesJSON, err := json.MarshalIndent(
		messages,
		"", "\t",
	)

	if err != nil {
		return nil, err
	}

	input := fmt.Sprintf(
		intentPrompt,
		string(schemaJSON),
		string(previousIntentJSON),
		string(messagesJSON),
	)

	input += fmt.Sprintf(
		"\n\nCurrent user request:\n%s",
		query,
	)

	response, err := c.Client.Responses.New(
		context.Background(),
		responses.ResponseNewParams{
			Model: openai.ChatModelGPT4_1Nano,
			Input: responses.ResponseNewParamsInputUnion{
				OfString: openai.String(input),
			},
		},
	)

	if err != nil {
		return nil, err
	}

	var result IntentResponse
	if err := json.Unmarshal(
		[]byte(response.OutputText()),
		&result,
	); err != nil {
		return nil, fmt.Errorf(
			"failed to parse intent response: %w\nresponse=%s",
			err,
			response.OutputText(),
		)
	}

	return &result, nil
}
