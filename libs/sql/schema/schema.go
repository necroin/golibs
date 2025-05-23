package sql_schema

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var (
	fieldsTypesMap = map[string]string{
		"null":      "null",
		"int":       "integer",
		"integer":   "integer",
		"float":     "float",
		"string":    "text",
		"text":      "text",
		"timestamp": "timestamp",
		"datetime":  "datetime",
	}
)

type TableField struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Nullable bool   `json:"nullable"`
}

type Table struct {
	Name       string       `json:"name"`
	Fields     []TableField `json:"fields"`
	PrimaryKey []string     `json:"primary_key"`
}

type SameTables struct {
	Current Table
	New     Table
}

type SameFields struct {
	Current TableField
	New     TableField
}

func Parse(schemaData []byte) ([]Table, error) {
	schema := []Table{}
	if err := json.Unmarshal(schemaData, &schema); err != nil {
		return nil, fmt.Errorf("failed parse schema data: %s", err)
	}
	return schema, nil
}

func Verify(tables []Table) error {
	for tableIndex, table := range tables {
		if table.Name == "" {
			return fmt.Errorf("table (%d) name is empty", tableIndex+1)
		}

		if len(table.Fields) == 0 {
			return fmt.Errorf(`missed "fields" section`)
		}

		fieldsMap := make(map[string]bool)
		for fieldIndex, field := range table.Fields {
			if field.Name == "" {
				return fmt.Errorf(`empty field (%d) name in "%s" table`, fieldIndex+1, table.Name)
			}
			_, ok := fieldsTypesMap[field.Type]
			if !ok {
				return fmt.Errorf(`unknown type for "%s" field in "%s" table`, field.Name, table.Name)
			}
			fieldsMap[field.Name] = true
		}

		if len(table.PrimaryKey) == 0 {
			return fmt.Errorf(`primary key for "%s" table is empty`, table.Name)
		}

		for _, fieldName := range table.PrimaryKey {
			_, ok := fieldsMap[fieldName]
			if !ok {
				return fmt.Errorf(`unknown field "%s" in primary key for "%s" table`, fieldName, table.Name)
			}
		}

	}

	return nil
}

func GetCurrentSchema(db *sql.DB) ([]byte, error) {
	result, err := db.Query("SELECT data FROM __Schema WHERE version = 'current'")
	if err != nil {
		return nil, fmt.Errorf("failed get current schema from migration table: %s", err)
	}
	defer result.Close()

	data := []byte{}
	if result.Next() {
		if err := result.Scan(&data); err != nil {
			return nil, fmt.Errorf("failed scan selected data from migration table: %s", err)
		}
	}
	if data == nil {
		return []byte{}, nil
	}

	return data, nil
}

func InitMetadata(db *sql.DB) error {
	table := Table{
		Name: "__Schema",
		Fields: []TableField{
			{
				Name:     "version",
				Type:     "string",
				Nullable: false,
			},
			{
				Name:     "data",
				Type:     "string",
				Nullable: true,
			},
		},
		PrimaryKey: []string{"version"},
	}

	if _, err := db.Exec(getCreateTableCommand(table)); err != nil {
		return fmt.Errorf("[SQL Schema] [Error] failed init migration table : %s", err)
	}

	db.Exec("INSERT INTO __Schema (version, data) VALUES('current','[]')")

	return nil
}

func findTableFieldsDifference(currentTable Table, newTable Table) ([]TableField, []TableField, []string) {
	newTableFieldsMap := make(map[string]TableField)
	currentTableFieldsMap := make(map[string]TableField)

	for _, field := range currentTable.Fields {
		currentTableFieldsMap[field.Name] = field
	}

	for _, field := range newTable.Fields {
		newTableFieldsMap[field.Name] = field
	}

	newFields := []TableField{}
	removedFields := []TableField{}
	sameNamedFields := []SameFields{}
	sameFields := []string{}

	for _, newField := range newTable.Fields {
		currentField, ok := currentTableFieldsMap[newField.Name]
		if ok {
			sameNamedFields = append(sameNamedFields, SameFields{Current: currentField, New: newField})
		}
		if !ok {
			newFields = append(newFields, newField)
		}
	}

	for _, field := range currentTable.Fields {
		_, ok := newTableFieldsMap[field.Name]
		if !ok {
			removedFields = append(removedFields, field)
		}
	}

	for _, fields := range sameNamedFields {
		if fields.New.Type != fields.Current.Type {
			newFields = append(newFields, fields.New)
			removedFields = append(removedFields, fields.Current)
		} else {
			sameFields = append(sameFields, fields.New.Name)
		}
	}

	return newFields, removedFields, sameFields
}

func findTablesDifference(currentSchema []Table, newSchema []Table) ([]Table, []Table, []SameTables) {
	newSchemaMap := make(map[string]Table)
	currentSchemaMap := make(map[string]Table)

	for _, newTable := range newSchema {
		newSchemaMap[newTable.Name] = newTable
	}

	for _, currentTable := range currentSchema {
		currentSchemaMap[currentTable.Name] = currentTable
	}

	newTables := []Table{}
	removedTables := []Table{}
	sameTables := []SameTables{}
	for _, newTable := range newSchema {
		currentTable, ok := currentSchemaMap[newTable.Name]
		if ok {
			sameTables = append(sameTables, SameTables{Current: currentTable, New: newTable})
		}
		if !ok {
			newTables = append(newTables, newTable)
		}
	}

	for _, currentTable := range currentSchema {
		_, ok := newSchemaMap[currentTable.Name]
		if !ok {
			removedTables = append(removedTables, currentTable)
		}
	}

	return newTables, removedTables, sameTables
}

func SchemaUpgrade(db *sql.DB, currentSchema []Table, newSchema []Table) error {
	newTables, removedTables, sameTables := findTablesDifference(currentSchema, newSchema)

	for _, table := range newTables {
		if _, err := db.Exec(getCreateTableCommand(table)); err != nil {
			return fmt.Errorf("[Schema Upgrade] failed create new table %s: %s", table.Name, err)
		}
	}

	for _, table := range removedTables {
		db.Exec(fmt.Sprintf("DROP TABLE %s", table.Name))
	}

	for _, tables := range sameTables {
		table := tables.New
		table.Name = "_new_" + table.Name

		newFields, removedFields, sameFields := findTableFieldsDifference(tables.Current, tables.New)
		tableIsChanged := len(newFields) != 0 || len(removedFields) != 0

		if tableIsChanged {
			if _, err := db.Exec(getCreateTableCommand(table)); err != nil {
				return fmt.Errorf(`[Schema Upgrade] failed create "%s" table: %s`, table.Name, err)
			}

			selectFields := strings.Join(sameFields, ", ")
			if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s", table.Name, selectFields, selectFields, tables.Current.Name)); err != nil {
				return fmt.Errorf(`[Schema Upgrade] failed copy data from "%s" table to "%s" table: %s`, tables.Current.Name, table.Name, err)
			}

			if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", tables.Current.Name)); err != nil {
				return fmt.Errorf(`[Schema Upgrade] failed delete "%s" table: %s`, tables.Current.Name, err)
			}

			if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %s RENAME TO %s", table.Name, tables.Current.Name)); err != nil {
				return fmt.Errorf(`[Schema Upgrade] failed rename "%s" table to "%s" table: %s`, table.Name, tables.Current.Name, err)
			}
		}
	}

	return nil
}

func Migration(db *sql.DB, currentSchema []Table, newSchema []Table) error {
	newSchemaData := &bytes.Buffer{}
	if err := json.NewEncoder(newSchemaData).Encode(&newSchema); err != nil {
		return fmt.Errorf("[Migration] failed encode new schema: %s", err)
	}

	if err := Transaction(db, func() error {
		if err := SchemaUpgrade(db, currentSchema, newSchema); err != nil {
			return err
		}

		if _, err := db.Exec("DELETE FROM __Schema WHERE version = 'current'"); err != nil {
			return fmt.Errorf("failed delete current schema from migration table: %s", err)
		}

		if _, err := db.Exec("INSERT INTO __Schema (version, data) VALUES ('current', $1)", newSchemaData.Bytes()); err != nil {
			return fmt.Errorf("failed insert new schema to migration table: %s", err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("[Migration] failed upgrade schema: %s", err)
	}

	return nil
}

func SetSchema(db *sql.DB, data io.Reader) error {
	if err := InitMetadata(db); err != nil {
		return fmt.Errorf("[SetSchema] failed init metadata: %s", err)
	}

	shcemaData, err := io.ReadAll(data)
	if err != nil {
		return fmt.Errorf("[SetSchema] failed read schema data: %s", err)
	}

	newSchema, err := Parse(shcemaData)
	if err != nil {
		return fmt.Errorf("[SetSchema] failed parse schema: %s", err)
	}

	if err := Verify(newSchema); err != nil {
		return fmt.Errorf("[SetSchema] failed verify schema: %s", err)
	}

	currentSchemaData, err := GetCurrentSchema(db)
	if err != nil {
		return fmt.Errorf("[SetSchema] failed get current schema: %s", err)
	}

	currentSchema, err := Parse(currentSchemaData)
	if err != nil {
		return fmt.Errorf("[SetSchema] failed parse current schema: %s", err)
	}

	if err := Migration(db, currentSchema, newSchema); err != nil {
		return fmt.Errorf("[SetSchema] failed migrate to new schema: %s", err)
	}

	return nil
}
