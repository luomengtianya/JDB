package utils

import (
	"encoding/json"
)

func ToJsonString(r interface{}) string {
	b, _ := json.Marshal(r)
	return string(b)
}

// IfNull 类似三元运算
func IfNull(target, v, other *string) string {

	if target == nil || *target == "" {
		return *other
	}

	return *v
}
