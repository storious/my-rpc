package serialization

import "encoding/json"

type JsonSerializer struct {
}

func (j JsonSerializer) Serialize(val any) ([]byte, error) {
	return json.Marshal(val)
}

func (j JsonSerializer) Deserialize(data []byte, ptr any) error {
	return json.Unmarshal(data, ptr)
}
