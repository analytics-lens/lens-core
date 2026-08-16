package model

import "time"

type Metric struct {
	Column       string   `json:"column"`
	Aggregations []string `json:"aggregations"`
	Type         string   `json:"type"`
}

type Dimension struct {
	Column string `json:"column"`
	Type   string `json:"type"`
}

type Schema struct {
	Metrics    map[string]Metric    `json:"metrics"`
	Dimensions map[string]Dimension `json:"dimensions"`
}

var SupportSchema = Schema{
	Metrics: map[string]Metric{
		"tickets": {
			Column: "tickets",
			Aggregations: []string{
				"sum",
				"avg",
			},
			Type: "integer",
		},

		"tickets_resolved": {
			Column: "tickets_resolved",
			Aggregations: []string{
				"sum",
				"avg",
			},
			Type: "integer",
		},

		"average_response_time": {
			Column: "average_response_time",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"first_response_time": {
			Column: "first_response_time",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"resolution_rate": {
			Column: "resolution_rate",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"automation_rate": {
			Column: "automation_rate",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"escalation_rate": {
			Column: "escalation_rate",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"csat": {
			Column: "csat",
			Aggregations: []string{
				"avg",
			},
			Type: "number",
		},

		"bot_conversions": {
			Column: "bot_conversions",
			Aggregations: []string{
				"sum",
				"avg",
			},
			Type: "integer",
		},

		"revenue_attributed_to_bot": {
			Column: "revenue_attributed_to_bot",
			Aggregations: []string{
				"sum",
				"avg",
			},
			Type: "number",
		},
	},

	Dimensions: map[string]Dimension{

		"date": {
			Column: "date",
			Type:   "date",
		},

		"channel": {
			Column: "channel",
			Type:   "string",
		},

		"agent": {
			Column: "agent",
			Type:   "string",
		},
	},
}

type SupportMetric struct {
	ID                        int64     `json:"id"`
	Date                      time.Time `json:"date"`
	Channel                   string    `json:"channel"`
	Agent                     string    `json:"agent"`
	Tickets                   int       `json:"tickets"`
	TicketsResolved           int       `json:"tickets_resolved"`
	AverageResponseTime       float64   `json:"average_response_time"`
	FirstResponseTime         float64   `json:"first_response_time"`
	ResolutionRate            float64   `json:"resolution_rate"`
	AutomationRate            float64   `json:"automation_rate"`
	EscalationRate            float64   `json:"escalation_rate"`
	CSAT                      float64   `json:"csat"`
	BotConversions            int       `json:"bot_conversions"`
	RevenueAttributedToBot    float64   `json:"revenue_attributed_to_bot"`
}