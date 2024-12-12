package msg_convert_protobuf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/necroin/golibs/libs/parser"
	"github.com/necroin/golibs/libs/tokenizer"
	"github.com/necroin/golibs/utils"
)

type TokenData struct {
	value string
	data  map[string]any
}

type Converter struct {
	tokenizer *tokenizer.Tokenizer
	parser    *parser.Parser[TokenData]
}

func NewConverter() *Converter {
	decoder := &Converter{}

	protobufTokenizer := tokenizer.NewTokenizer(
		tokenizer.NewToken("WORD", `[a-zA-Z_0-9][a-zA-Z0-9_\-]*`),
		tokenizer.NewToken("QUOTE", `\"`),
		tokenizer.NewToken("OBJECT_OPEN_BRACKET", `\<`),
		tokenizer.NewToken("OBJECT_CLOSE_BRACKET", `\>`),
		tokenizer.NewToken("COLON", `\:`),
	)
	decoder.tokenizer = protobufTokenizer

	protobufParser := parser.NewParser[TokenData]()
	decoder.parser = protobufParser

	protobufParser.AddRule(parser.NewRule("VALUE", "WORD", func(tokens []parser.Token[TokenData]) TokenData {
		return tokens[0].Value()
	}))

	protobufParser.AddRule(parser.NewRule("VALUE", "NUMBER", func(tokens []parser.Token[TokenData]) TokenData {
		return tokens[0].Value()
	}))

	protobufParser.AddRule(parser.NewRule("VALUE", "QUOTE QUOTE", func(tokens []parser.Token[TokenData]) TokenData {
		return TokenData{
			value: "",
		}
	}))

	protobufParser.AddRule(parser.NewRule("VALUE", "QUOTE VALUE QUOTE", func(tokens []parser.Token[TokenData]) TokenData {
		return tokens[1].Value()
	}))

	protobufParser.AddRule(parser.NewRule("VALUE", "OBJECT_OPEN_BRACKET KEY_VALUE_LIST OBJECT_CLOSE_BRACKET", func(tokens []parser.Token[TokenData]) TokenData {
		return tokens[1].Value()
	}))

	protobufParser.AddRule(parser.NewRule("KEY_VALUE", "VALUE COLON VALUE", func(tokens []parser.Token[TokenData]) TokenData {
		keyToken := tokens[0].Value()
		valueToken := tokens[2].Value()

		if valueToken.data != nil {
			return TokenData{data: map[string]any{keyToken.value: valueToken.data}}
		}

		return TokenData{data: map[string]any{keyToken.value: valueToken.value}}
	}))

	protobufParser.AddRule(parser.NewRule("KEY_VALUE_LIST", "KEY_VALUE", func(tokens []parser.Token[TokenData]) TokenData {
		token := tokens[0].Value()

		return TokenData{
			data: token.data,
		}
	}))

	protobufParser.AddRule(parser.NewRule("KEY_VALUE_LIST", "KEY_VALUE_LIST KEY_VALUE", func(tokens []parser.Token[TokenData]) TokenData {
		listToken := tokens[0].Value()
		kvToken := tokens[1].Value()

		for key, value := range kvToken.data {
			listTokenValue, ok := listToken.data[key]
			if !ok {
				listToken.data[key] = value
				continue
			}

			if utils.IsSlice(reflect.ValueOf(listTokenValue)) {
				listTokenValueSlice := listTokenValue.([]any)
				listTokenValueSlice = append(listTokenValueSlice, value)
				listToken.data[key] = listTokenValueSlice
				continue
			}

			valueSlice := []any{
				listTokenValue,
				value,
			}
			listToken.data[key] = valueSlice
		}

		return listToken
	}))

	return decoder
}

func (converter *Converter) ToJson(data []byte, out io.Writer) error {
	tokens, err := converter.tokenizer.Parse(data)
	if err != nil {
		return fmt.Errorf("failed tokenize data: %s", err)
	}

	parserTokens := []parser.Token[TokenData]{}
	for _, token := range tokens {
		parserTokens = append(parserTokens, parser.NewParserToken(token.Name(), TokenData{value: token.Value()}))
	}

	resultToken, err := converter.parser.Parse(parser.ParseOptions{}, parserTokens...)
	if err != nil {
		return fmt.Errorf("failed parse tokens: %s", err)
	}

	result, _ := json.Marshal(resultToken.Value().data)

	buffer := &bytes.Buffer{}
	json.Indent(buffer, result, "", "\t")
	out.Write(buffer.Bytes())
	return nil
}
