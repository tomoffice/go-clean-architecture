// sensitive.go
package logging

import "strings"

type SensitiveKeys struct {
	Keys map[string]struct{}
}

func NewSensitiveKeys(keys []string) *SensitiveKeys {
	m := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		m[strings.ToLower(k)] = struct{}{}
	}
	return &SensitiveKeys{Keys: m}
}
func (s *SensitiveKeys) isSensitive(key string) bool {
	_, ok := s.Keys[strings.ToLower(key)]
	return ok
}

// ★ 遞迴遮罩：支援 map / slice（JSON 物件與陣列）、支援巢狀
func (s *SensitiveKeys) Mask(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, vv := range x {
			if s.isSensitive(k) {
				out[k] = "******"
				continue
			}
			out[k] = s.Mask(vv)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i := range x {
			out[i] = s.Mask(x[i])
		}
		return out
	default:
		return v
	}
}
