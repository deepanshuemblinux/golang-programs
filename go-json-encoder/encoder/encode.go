package encoder

import (
	"fmt"
	"io"
	"reflect"
	"strings"
)

type JSONEncoder struct {
	writer io.WriteCloser
}

func NewEncoder(writer io.WriteCloser) *JSONEncoder {
	return &JSONEncoder{
		writer: writer,
	}
}

func (enc *JSONEncoder) writeVal(val reflect.Value, builder *strings.Builder) {
	switch val.Type().Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		builder.WriteString(fmt.Sprintf("%d", val.Int()))
	case reflect.String:
		builder.WriteString(fmt.Sprintf("\"%s\"", val.String()))
	case reflect.Float32, reflect.Float64:
		builder.WriteString(fmt.Sprintf("%f", val.Float()))
	case reflect.Map:

		enc.handleMap(val.Interface(), builder)
	case reflect.Slice:

		enc.handleSlice(val.Interface(), builder)
	case reflect.Struct:
		enc.handleStruct(val.Interface(), builder)

	}
}
func (enc *JSONEncoder) handleStruct(data interface{}, builder *strings.Builder) error {
	builder.WriteString("{")
	dVal := reflect.ValueOf(data)
	numF := dVal.NumField()
	for i := 0; i < numF; i++ {
		fld := dVal.Field(i)
		if name, found := dVal.Type().Field(i).Tag.Lookup("myjson"); found {
			builder.WriteString(fmt.Sprintf("\"%s\": ", name))
		} else {
			builder.WriteString(fmt.Sprintf("\"%s\": ", dVal.Type().Field(i).Name))
		}
		enc.writeVal(fld, builder)
		if i < numF-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("}")
	return nil
}
func (enc *JSONEncoder) handleSlice(data interface{}, builder *strings.Builder) error {
	builder.WriteString("[")
	dVal := reflect.ValueOf(data)
	len := dVal.Len()
	for i := 0; i < len; i++ {
		val := dVal.Index(i)
		enc.writeVal(val, builder)
		if i < len-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("]")
	return nil
}
func (enc *JSONEncoder) handleMap(data interface{}, builder *strings.Builder) error {
	builder.WriteString("{")
	dVal := reflect.ValueOf(data)
	keys := dVal.MapKeys()
	mLen := len(keys)
	for idx, k := range keys {
		enc.writeVal(k, builder)
		builder.WriteString(": ")
		val := dVal.MapIndex(k)
		enc.writeVal(val, builder)
		if idx < mLen-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("}")
	return nil
}
func (enc *JSONEncoder) encode(data interface{}, builder *strings.Builder) error {
	dType := reflect.TypeOf(data)
	kind := dType.Kind()
	switch kind {
	case reflect.Map:
		enc.handleMap(data, builder)
	case reflect.Struct:
		enc.handleStruct(data, builder)
	default:
		return fmt.Errorf("invalid data passed for encoding")
	}
	return nil
}
func (enc *JSONEncoder) Encode(data interface{}) error {
	builder := strings.Builder{}
	err := enc.encode(data, &builder)
	if err != nil {
		return err
	}
	enc.writer.Write([]byte(builder.String()))
	enc.writer.Close()
	return nil
}
