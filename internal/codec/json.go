package codec

import (
	"encoding/json"

	"google.golang.org/grpc/encoding"
)

const Name = "json"

func init() {
	encoding.RegisterCodec(&JSONCodec{})
}

type JSONCodec struct{}

func (j *JSONCodec) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j *JSONCodec) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func (j *JSONCodec) Name() string {
	return Name
}
