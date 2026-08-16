package services

type QueryField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func InferFields(rows []map[string]any) []QueryField {
	if len(rows) == 0 {
		return []QueryField{}
	}

	fields := make(
		[]QueryField,
		0,
		len(rows[0]),
	)

	for key := range rows[0] {

		field := QueryField{
			Name: key,
			Type: "unknown",
		}

		for _, row := range rows {

			value := row[key]
			if value == nil {
				continue
			}

			switch value.(type) {
			case int,
				int8,
				int16,
				int32,
				int64,
				uint,
				uint8,
				uint16,
				uint32,
				uint64:
				field.Type = "integer"

			case float32,
				float64:
				field.Type = "number"

			case string:
				field.Type = "string"
			}

			if field.Type != "unknown" {
				break
			}
		}

		fields = append(
			fields,
			field,
		)
	}

	return fields
}
