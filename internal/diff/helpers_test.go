package diff

import (
	"encoding/json"
	"strings"
)

// parseJSON is a test helper that parses JSON with UseNumber for precision.
func parseJSON(s string) interface{} {
	var v interface{}
	decoder := json.NewDecoder(strings.NewReader(s))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		panic("parseJSON: " + err.Error())
	}
	return v
}

// newCtx creates a test context with default options.
func newCtx() *Context {
	return &Context{Opts: Options{}, Dispatcher: Dispatch}
}

// newCtxWithKey creates a test context with a key field.
func newCtxWithKey(key string) *Context {
	return &Context{Opts: Options{KeyField: key}, Dispatcher: Dispatch}
}
