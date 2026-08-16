package services

import (
	"fmt"
	"slices"
	"strings"

	model "core/internal/model"
)

type Query struct {
	SQL    string
	Params map[string]any
}

func validateIntent(intent QueryIntent) error {

	if intent.Metric == "" {
		return fmt.Errorf("missing metric")
	}

	metric, exists := model.SupportSchema.Metrics[intent.Metric]

	if !exists {
		return fmt.Errorf(
			"unknown metric: %s",
			intent.Metric,
		)
	}

	if intent.Aggregation == "" {
		return fmt.Errorf(
			"missing aggregation",
		)
	}

	aggregationAllowed := slices.Contains(
		metric.Aggregations,
		intent.Aggregation,
	)

	if !aggregationAllowed {
		return fmt.Errorf(
			"invalid aggregation %q for %s",
			intent.Aggregation,
			intent.Metric,
		)
	}

	if intent.Dimension != nil &&
		*intent.Dimension != "" {

		if _, exists :=
			model.SupportSchema.Dimensions[*intent.Dimension]; !exists {

			return fmt.Errorf(
				"unknown dimension: %s",
				*intent.Dimension,
			)
		}
	}

	for _, filter := range intent.Filters {

		if _, exists :=
			model.SupportSchema.Dimensions[filter.Dimension]; !exists {

			return fmt.Errorf(
				"unknown filter dimension: %s",
				filter.Dimension,
			)
		}

		if filter.Operator != "eq" {
			return fmt.Errorf(
				"only eq filters are supported",
			)
		}

		if filter.Value == "" {
			return fmt.Errorf(
				"empty filter value",
			)
		}
	}

	return nil
}

func BuildSQL(intent QueryIntent) (*Query, error) {

	if err := validateIntent(intent); err != nil {
		return nil, err
	}

	metric :=
		model.SupportSchema.Metrics[intent.Metric]

	selects := []string{}

	var dimension *model.Dimension

	if intent.Dimension != nil &&
		*intent.Dimension != "" {

		d := model.SupportSchema.Dimensions[*intent.Dimension]

		dimension = &d

		selects = append(
			selects,
			fmt.Sprintf(
				"%s AS %s",
				d.Column,
				*intent.Dimension,
			),
		)
	}

	selects = append(
		selects,
		fmt.Sprintf(
			"%s(%s) AS %s",
			strings.ToUpper(intent.Aggregation),
			metric.Column,
			intent.Metric,
		),
	)

	sql := fmt.Sprintf(
		`SELECT %s FROM support_metrics`,
		strings.Join(selects, ", "),
	)

	conditions := []string{}
	params := map[string]any{}

	switch intent.TimeRange.Type {
	case "last_7_days":
		conditions = append(
			conditions,
			`date >= (
					SELECT MAX(date)
					FROM support_metrics
				) - INTERVAL '6 days'`,
		)

	case "last_30_days":
		conditions = append(
			conditions,
			`date >= (
					SELECT MAX(date)
					FROM support_metrics
				) - INTERVAL '29 days'`,
		)

	case "all", "":
	default:
		return nil, fmt.Errorf(
			"unsupported time range: %s",
			intent.TimeRange.Type,
		)
	}

	for index, filter := range intent.Filters {
		dimension := model.SupportSchema.Dimensions[filter.Dimension]

		parameter := fmt.Sprintf(
			"filter_%d",
			index,
		)

		conditions = append(
			conditions,
			fmt.Sprintf(
				"%s = @%s",
				dimension.Column,
				parameter,
			),
		)

		params[parameter] = filter.Value
	}

	if len(conditions) > 0 {
		sql += ` WHERE ` + strings.Join(
			conditions,
			" AND ",
		)
	}

	if dimension != nil {
		sql += fmt.Sprintf(
			` GROUP BY %s ORDER BY %s DESC`,
			dimension.Column,
			intent.Metric,
		)
	}

	return &Query{
		SQL:    sql,
		Params: params,
	}, nil
}
