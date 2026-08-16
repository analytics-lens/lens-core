package services

type QueryIntent struct {
	Metric      string            `json:"metric"`
	Aggregation string            `json:"aggregation"`
	Dimension   *string           `json:"dimension"`
	TimeRange   TimeRange         `json:"timeRange"`
	Filters     []QueryFilter     `json:"filters"`
}

type TimeRange struct {
	Type  string  `json:"type"`
	Field *string `json:"field"`
}

type QueryFilter struct {
	Dimension string `json:"dimension"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
}