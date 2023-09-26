package sqlschema

func MapSchema(schema []Table) map[string]Table {
	result := map[string]Table{}
	for _, table := range schema {
		result[table.Name] = table
	}
	return result
}
