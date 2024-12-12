package msg_convert_protobuf_tests

import (
	"os"
	"testing"

	msg_convert_protobuf "github.com/necroin/golibs/libs/msg_convert/protobuf"
)

func TestDecode(t *testing.T) {
	message := "field_1:\"field_1_value\" field_2:\"12345678\" field_3:12345678 field_4:<field_4_field_1:\"field_4_field_1_value\" field_4_field_2:\"12345678\" field_4_field_3:12345678 field_4_field_4:true> field_4:<field_4_field_1:\"field_4_field_1_value_2\" field_4_field_2:\"87654321\" field_4_field_3:87654321 field_4_field_4:false>"

	converter := msg_convert_protobuf.NewConverter()
	if err := converter.ToJson([]byte(message), os.Stdout); err != nil {
		t.Fatal(err)
	}
}
