package kafka_test

import (
	"testing"

	"github.com/bhmt/tittlemanscrest/kafka"
	"github.com/stretchr/testify/assert"
)

type TestEvent struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func TestJSONSerializerIsOk(t *testing.T) {
	serializer := kafka.JsonSerializer[TestEvent]{}
	original := TestEvent{Name: "Test", Value: 100}

	// 1. Test Serialization
	data, err := serializer.Serialize(original)
	assert.NoError(t, err)
	assert.True(t, len(data) > 0)

	// Verify the output data structure manually
	expectedJSON := `{"name":"Test","value":100}`
	assert.JSONEq(t, expectedJSON, string(data))

	// 2. Test Deserialization
	restored, err := serializer.Deserialize(data)
	assert.NoError(t, err)
	assert.Equal(t, original, restored)
}
