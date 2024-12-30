package serialization

type Serializer interface {
	Serialize(any) ([]byte, error)
	Deserialize(data []byte, ptr any) error
}
