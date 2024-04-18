package rstruct_tests

type CommonExtendStruct struct {
	FirstField  string `json:"first_field"`
	SecondField int    `json:"second_field"`
	ThirdField  bool   `json:"third_field"`
}

type PointerExtendStruct struct {
	FirstField  *string `json:"first_field"`
	SecondField *int    `json:"second_field"`
	ThirdField  *bool   `json:"third_field"`
}

type NestedExtendStruct struct {
	NestedFirstField struct {
		FirstField string `json:"first_field"`
	} `json:"nested_first_field"`
	NestedSecondField struct {
		SecondField int `json:"second_field"`
	} `json:"nested_second_field"`
	NestedThirdField *struct {
		ThirdField bool `json:"third_field"`
	} `json:"nested_third_field"`
	NotNestedField string `json:"not_nested_field"`
}

type UnExportedFieldsStruct struct {
	unExportedField any
	ExportedField   string `json:"exported_field"`
}

type SimpleNestedStruct struct {
	FirstField  string `json:"first_field"`
	SecondField int    `json:"second_field"`
}

type IgnoreNestedStruct struct {
	NestedFirstField struct {
		FirstField string `json:"first_field"`
	} `json:"nested_first_field"`
	NestedSecondField SimpleNestedStruct  `json:"nested_second_field"`
	NestedThirdField  *SimpleNestedStruct `json:"nested_third_field"`
}
