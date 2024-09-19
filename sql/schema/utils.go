package sqlschema

import (
	"database/sql"
	"fmt"
	"strings"
)

func MapSchema(schema []Table) map[string]Table {
	result := map[string]Table{}
	for _, table := range schema {
		result[table.Name] = table
	}
	return result
}

func getCreateTableCommand(table Table) string {
	fields := []string{}
	for _, field := range table.Fields {
		tableFieldType := fieldsTypesMap[field.Type]
		tableField := fmt.Sprintf("%s %s NOT NULL", field.Name, tableFieldType)
		if field.Nullable {
			tableField = fmt.Sprintf("%s %s NULL", field.Name, tableFieldType)
		}
		fields = append(fields, tableField)
	}

	return fmt.Sprintf(
		"CREATE TABLE IF NOT EXISTS %s(%s, PRIMARY KEY(%s))",
		table.Name,
		strings.Join(fields, ", "),
		strings.Join(table.PrimaryKey, ", "),
	)
}

func Transaction(db *sql.DB, handler func() error) error {
	if _, err := db.Exec("BEGIN TRANSACTION"); err != nil {
		return fmt.Errorf("[Transaction] failed start transaction: %s", err)
	}

	if err := handler(); err != nil {
		if _, err := db.Exec("ROLLBACK TRANSACTION"); err != nil {
			return fmt.Errorf("[Transaction] failed rollback transaction: %s", err)
		}
		return err
	}

	if _, err := db.Exec("COMMIT TRANSACTION"); err != nil {
		return fmt.Errorf("[Transaction] failed commit transaction: %s", err)
	}

	return nil
}
