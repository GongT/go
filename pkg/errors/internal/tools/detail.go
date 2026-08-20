package tools

import (
	"fmt"

	"github.com/gongt/go/pkg/i18n/type_name"
)

func ApplyPairsToMap(m map[string]any, pairs []any) {
	if len(pairs)%2 != 0 {
		panic(fmt.Sprintf("ApplyPairsToMap: 参数数量必须为偶数，但得到%d个", len(pairs)))
	}

	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			panic(fmt.Sprintf("ApplyPairsToMap: 键必须是字符串，而不能是%s", type_name.TranslateInterfaceType(pairs[i])))
		}
		value := pairs[i+1]
		m[key] = value
	}
}
