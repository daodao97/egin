package json

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

var data = `
{
	"a": {
		"b": false,
		"c": 1,
		"d": "true"
	}
}
`

func TestJsonObj_GetBool(t *testing.T) {
	obj := NewJsonObj(data)
	// assert.Equal(t, obj.GetBool("a.b", false), true)
	assert.Equal(t, obj.GetBool("a.d", false), true)
}
