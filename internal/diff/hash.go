package diff

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// CanonicalJSON converts a JSON value to a canonical string representation
// where object keys are sorted alphabetically. This ensures that
// {"a":1,"b":2} and {"b":2,"a":1} produce the same hash.
func CanonicalJSON(v interface{}) string {
	var sb strings.Builder
	writeCanonical(&sb, v)
	return sb.String()
}

func writeCanonical(sb *strings.Builder, v interface{}) {
	if v == nil {
		sb.WriteString("null")
		return
	}

	switch val := v.(type) {
	case map[string]interface{}:
		sb.WriteByte('{')
		keys := make([]string, 0, len(val))
		for k := range val {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for i, k := range keys {
			if i > 0 {
				sb.WriteByte(',')
			}
			keyJSON, _ := json.Marshal(k)
			sb.Write(keyJSON)
			sb.WriteByte(':')
			writeCanonical(sb, val[k])
		}
		sb.WriteByte('}')

	case []interface{}:
		sb.WriteByte('[')
		for i, elem := range val {
			if i > 0 {
				sb.WriteByte(',')
			}
			writeCanonical(sb, elem)
		}
		sb.WriteByte(']')

	case json.Number:
		sb.WriteString(val.String())

	case string:
		b, _ := json.Marshal(val)
		sb.Write(b)

	case bool:
		if val {
			sb.WriteString("true")
		} else {
			sb.WriteString("false")
		}

	default:
		sb.WriteString(fmt.Sprintf("%v", val))
	}
}

// HashJSON returns a SHA256 hash of the canonical JSON representation
func HashJSON(v interface{}) string {
	canonical := CanonicalJSON(v)
	h := sha256.Sum256([]byte(canonical))
	return fmt.Sprintf("%x", h)
}
