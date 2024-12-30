package serialization

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"myRPC/rpc-core/model"
	"reflect"
)

var Space = [...]byte{19, 23, 29, 31, 51}

func IntToBytes(i int) ([]byte, error) {
	x := int64(i)
	bytesBuffer := bytes.NewBuffer([]byte{})
	err := binary.Write(bytesBuffer, binary.BigEndian, x)
	if err != nil {
		return nil, errors.New("binary.Write failed")
	}
	return bytesBuffer.Bytes(), nil
}

func BytesToInt(b []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	err := binary.Read(bytesBuffer, binary.BigEndian, &x)
	if err != nil {
		return 0, errors.New("binary.Read failed")
	}
	return int(x), nil
}

func MarshalArguments(args []interface{}) ([]byte, error) {
	types := make([]byte, 0, len(args))
	lens := make([]int, 0, len(args))
	buf := make([]byte, 0, 256)
	buffer := bytes.NewBuffer(buf)
	for _, arg := range args {
		switch v := arg.(type) {
		case int:
			types = append(types, model.TypeInt)
			lens = append(lens, 8)
			err := binary.Write(buffer, binary.BigEndian, int64(v))
			if err != nil {
				return nil, err
			}
		case string:
			types = append(types, model.TypeString)
			lens = append(lens, len(v))
			buffer.WriteString(v)
		case float32:
			err := binary.Write(buffer, binary.BigEndian, v)
			if err != nil {
				return nil, err
			}
			types = append(types, model.TypeFloat32)
			lens = append(lens, 4)
		case float64:
			err := binary.Write(buffer, binary.BigEndian, v)
			if err != nil {
				return nil, err
			}
			types = append(types, model.TypeFloat64)
			lens = append(lens, 8)
		case bool:
			err := binary.Write(buffer, binary.BigEndian, v)
			if err != nil {
				return nil, err
			}
			types = append(types, model.TypeBool)
			lens = append(lens, 1)
		default:
			return nil, errors.New("unsupported type")
		}
	}

	result := make([]byte, 0, len(Space)+8+len(types)+8*len(types)+buffer.Len())
	writer := bytes.NewBuffer(result)
	writer.Write(Space[:])           // array to slice
	res, _ := IntToBytes(len(types)) // len(types)
	writer.Write(res)                // args types

	for _, t := range types {
		writer.WriteByte(t)
	}

	for _, l := range lens {
		res, _ = IntToBytes(l)
		writer.Write(res)
	}
	writer.Write(buffer.Bytes())
	return writer.Bytes(), nil
}

func UnmarshalArguments(stream []byte) ([]any, error) {
	if len(stream) < len(Space)+8 {
		return nil, errors.New("invalid stream: too short")
	}
	if !bytes.Equal(stream[:len(Space)], Space[:]) {
		return nil, errors.New("invalid stream: magic number verification failed")
	}
	pos := len(Space)
	n, _ := BytesToInt(stream[pos : pos+8])
	pos += 8
	if n <= 0 {
		return nil, nil
	}
	if len(stream) < len(Space)+8+n+8*n {
		return nil, errors.New("invalid stream: too short")
	}
	types := make([]byte, 0, n)
	for i := 0; i < n; i++ {
		types = append(types, stream[pos])
		pos++
	}
	totalLength := 0
	lens := make([]int, 0, n) // lens of each arg
	for i := 0; i < n; i++ {
		l, _ := BytesToInt(stream[pos : pos+8])
		pos += 8
		lens = append(lens, l)
		totalLength += l
	}

	if len(stream[pos:]) < totalLength {
		return nil, errors.New("invalid stream: too short")
	}

	args := make([]any, 0, n) // content of each arg
	for i := 0; i < n; i++ {
		t := types[i]
		l := lens[i]
		bytesBuffer := bytes.NewBuffer(stream[pos : pos+l])

		var arg any
		switch t {
		case model.TypeBool:
			var v bool
			err := binary.Read(bytesBuffer, binary.BigEndian, &v)
			if err != nil {
				return nil, err
			}
			arg = v
		case model.TypeFloat32:
			var v float32
			err := binary.Read(bytesBuffer, binary.BigEndian, &v)
			if err != nil {
				return nil, err
			}
			arg = v
		case model.TypeFloat64:
			var v float64
			err := binary.Read(bytesBuffer, binary.BigEndian, &v)
			if err != nil {
				return nil, err
			}
			arg = v

		case model.TypeInt:
			x, _ := BytesToInt(stream[pos : pos+l])
			arg = x
		case model.TypeString:
			arg = string(stream[pos : pos+l])
		default:
			return nil, errors.New("unsupported type")
		}
		args = append(args, arg)
		pos += l
	}
	return args, nil
}

type MySerializer struct {
}

func (m MySerializer) Serialize(object any) ([]byte, error) {
	args := make([]any, 0, 10)
	typ := reflect.TypeOf(object)
	val := reflect.ValueOf(object)
	for i := 0; i < typ.NumField(); i++ {
		if !typ.Field(i).IsExported() {
			continue
		}
		args = append(args, val.Field(i).Interface())
	}
	return MarshalArguments(args)

}
func (m MySerializer) Deserialize(stream []byte, objet any) error {
	typ := reflect.TypeOf(objet)
	val := reflect.ValueOf(objet)
	if typ.Kind() != reflect.Ptr {
		return errors.New("objet must be a pointer")
	}

	typ = typ.Elem() // parse ptr
	val = val.Elem()

	if typ.Kind() != reflect.Struct {
		return errors.New("objet must be a struct")
	}
	args, err := UnmarshalArguments(stream)
	if err != nil {
		return err
	}
	j := 0
	for i := 0; i < typ.NumField(); i++ {
		if !typ.Field(i).IsExported() {
			continue
		}
		arg := args[j]
		j++
		switch typ.Field(i).Type.Kind() {
		case reflect.Int:
			if v, ok := arg.(int); ok {
				val.Field(i).SetInt(int64(v))
			} else {
				return fmt.Errorf("type missmatch: expect int, got %T", arg)
			}
		case reflect.Float32:
			if v, ok := arg.(float32); ok {
				val.Field(i).SetFloat(float64(v))
			} else {
				return fmt.Errorf("type missmatch: expect float32, got %T", arg)
			}
		case reflect.Float64:
			if v, ok := arg.(float64); ok {
				val.Field(i).SetFloat(v)
			} else {
				return fmt.Errorf("type missmatch: expect float64, got %T", arg)
			}
		case reflect.String:
			if v, ok := arg.(string); ok {
				val.Field(i).SetString(v)
			} else {
				return fmt.Errorf("type missmatch: expect string, got %T", arg)
			}
		case reflect.Bool:
			if v, ok := arg.(bool); ok {
				val.Field(i).SetBool(v)
			} else {
				return fmt.Errorf("type missmatch: expect bool, got %T", arg)
			}
		default:
			return fmt.Errorf("unsupported type: %s", typ.Field(i).Type.Kind())
		}
	}
	return nil
}
