package iterator

import (
	"log"
	"strconv"
)

func GetCode(err error) (int, bool) {
	for details := range IterEveryDetail(err) {
		if code, ok := details["code"]; ok {
			if codeInt, ok := convertErrorCode(code); ok {
				// log.Println("找到错误码:", codeInt)
				return codeInt, true
			}
		}
		// 没有code字段或不可用，继续查找
	}

	// log.Println("未找到错误码")
	return -1, false
}

// 将所有整数类型转换为int
func int_conv(value any) int {
	switch v := value.(type) {
	case int:
		return v
	case int8:
		return int(v)
	case int16:
		return int(v)
	case int32:
		return int(v)
	case int64:
		return int(v)
	case uint:
		return int(v)
	case uint8:
		return int(v)
	case uint16:
		return int(v)
	case uint32:
		return int(v)
	case uint64:
		return int(v)
	}
	return 0
}

func convertErrorCode(value any) (codeInt int, ok bool) {
	switch code := value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		codeInt = int_conv(code)
		ok = true
	case string:
		// 尝试将字符串转换为整数
		if code == "" {
			return -1, false
		}
		var parsedCode int
		parsedCode, err := strconv.Atoi(code)
		if err == nil {
			codeInt = parsedCode
			ok = true
		} else {
			log.Printf(`无法将错误码字符串"%s"转换为整数: %v\n`, code, err)
			ok = false
		}
	default:
		log.Printf("未知的错误码类型: %#v\n", code)
		ok = false
	}
	if ok == false && codeInt == 0 {
		codeInt = -1
	}
	return codeInt, ok
}

func GetReason(err error) (reason error, found bool) {
	for details := range IterEveryDetail(err) {
		if reason, ok := details["reason"]; ok {
			if reasonErr, ok := reason.(error); ok {
				return reasonErr, true
			}
		}
		// 没有reason字段或不可用，继续查找
	}

	return nil, false
}
