package kafka

import (
	"encoding/json"

	"google.golang.org/protobuf/proto"
)

type Serializer[T any] interface {
	Serialize(T) ([]byte, error)
	Deserialize([]byte) (T, error)
}

type JsonSerializer[T any] struct{}

func (j JsonSerializer[T]) Serialize(data T) ([]byte, error) {
	return json.Marshal(data)
}

func (j JsonSerializer[T]) Deserialize(data []byte) (T, error) {
	var target T
	err := json.Unmarshal(data, &target)
	return target, err
}

type ProtoSerializer[T proto.Message] struct{}

func (j ProtoSerializer[T]) Serialize(data T) ([]byte, error) {
	return proto.Marshal(data)
}

func (j ProtoSerializer[T]) Deserialize(data []byte) (T, error) {
	var target T
	err := proto.Unmarshal(data, target)
	return target, err
}
