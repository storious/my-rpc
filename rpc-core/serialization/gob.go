package serialization

import (
	"bytes"
	"encoding/gob"
)

type Gob struct {
}

func (g Gob) Serialize(a any) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(a); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (g Gob) Deserialize(data []byte, ptr any) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	return decoder.Decode(ptr)
}
