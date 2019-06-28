package interfaces

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestInterfaces(t *testing.T) {
	t.Run("Can unmarshal LogFields", func(*testing.T) {
		t.Run("with strings as values", func(*testing.T) {
			var fields LogFields
			err := json.Unmarshal(([]byte)("[\"a\", \"b\"]"), &fields)
			if err != nil {
				t.Error("Unmarshaling failed", err)

			}
			if !reflect.DeepEqual(LogFields{"a", "b"}, fields) {
				t.Error("Got wrong values", fields)
			}

		})
		t.Run("with numbers as values", func(*testing.T) {
			var fields LogFields
			err := json.Unmarshal(([]byte)("[1, 2]"), &fields)
			if err != nil {
				t.Error("Unmarshaling failed", err)

			}
			if !reflect.DeepEqual(LogFields{"1", "2"}, fields) {
				t.Error("Got wrong values", fields)
			}

		})
	})

	t.Run("Can unmarshal LogGeneric", func(*testing.T) {
		t.Run("with strings as values", func(*testing.T) {
			var fields LogGeneric
			err := json.Unmarshal(([]byte)("{\"key1\":\"value1\", \"key2\":\"value2\"}"), &fields)
			if err != nil {
				t.Error("Unmarshaling failed", err)

			}
			if !reflect.DeepEqual(LogGeneric{"key1": "value1", "key2": "value2"}, fields) {
				t.Error("Got wrong values", fields)
			}
		})
		t.Run("with numbers as values", func(*testing.T) {
			var fields LogGeneric
			err := json.Unmarshal(([]byte)("{\"key1\":1, \"key2\":2}"), &fields)
			if err != nil {
				t.Error("Unmarshaling failed", err)
			}
			if !reflect.DeepEqual(LogGeneric{"key1": "1", "key2": "2"}, fields) {
				t.Error("Got wrong values", fields)
			}
		})
	})
}
