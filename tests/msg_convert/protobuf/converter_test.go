package msg_convert_protobuf_tests

import (
	"bytes"
	"strings"
	"testing"

	msg_convert_protobuf "github.com/necroin/golibs/libs/msg_convert/protobuf"
)

func TestJson(t *testing.T) {
	message := "field_1:\"field_1_value\" field_2:\"12345678\" field_3:12345678 field_4:<field_4_field_1:\"field_4_field_1_value\" field_4_field_2:\"12345678\" field_4_field_3:12345678 field_4_field_4:true> field_4:<field_4_field_1:\"field_4_field_1_value_2\" field_4_field_2:\"87654321\" field_4_field_3:87654321 field_4_field_4:false>"

	expected := `
		{
    	    "field_1": "field_1_value",
    	    "field_2": "12345678",
    	    "field_3": "12345678",
    	    "field_4": [
    	            {
    	                    "field_4_field_1": "field_4_field_1_value",
    	                    "field_4_field_2": "12345678",
    	                    "field_4_field_3": "12345678",
    	                    "field_4_field_4": "true"
    	            },
    	            {
    	                    "field_4_field_1": "field_4_field_1_value_2",
    	                    "field_4_field_2": "87654321",
    	                    "field_4_field_3": "87654321",
    	                    "field_4_field_4": "false"
    	            }
    	    ]
		}
	`

	expected = strings.ReplaceAll(expected, "\t", "")
	expected = strings.ReplaceAll(expected, "\n", "")
	expected = strings.ReplaceAll(expected, " ", "")

	converter := msg_convert_protobuf.NewConverter()
	resultBuffer := &bytes.Buffer{}
	if err := converter.ToJson([]byte(message), resultBuffer); err != nil {
		t.Fatal(err)
	}

	result := resultBuffer.String()
	result = strings.ReplaceAll(result, "\t", "")
	result = strings.ReplaceAll(result, "\n", "")
	result = strings.ReplaceAll(result, " ", "")

	if result != expected {
		t.Fatalf("messages not equal, %s != %s", result, expected)
	}
}
