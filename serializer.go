package fastlocalcache

import "encoding/json"

type Serializer interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(data []byte, v any) error
}

type JSONSerializer struct {
}

func (j JSONSerializer) Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func (j JSONSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}
