package utils

import "github.com/go-park-mail-ru/2022_1_CJ/internal/model/core"

func GetLimitArray(array *[]string, limit, page int64) ([]string, int64, int64) {
	total := int64(len(*array))

	if limit != -1 && limit <= total {
		start := total - limit*(page)
		end := total - limit*(page-1)
		if end < 0 {
			end = 0
		}
		if start < 0 {
			start = 0
		}
		limitArray := (*array)[start:end]

		return reverseString(limitArray), total, total/limit + IsLarge(total%limit > 0)
	} else {
		return reverseString(*array), total, 1
	}
}

func GetLimitMessage(array *[]core.Message, limit, page int64) ([]core.Message, int64, int64) {
	total := int64(len(*array))

	if limit != -1 && limit <= total {
		start := total - limit*(page)
		end := total - limit*(page-1)
		if end < 0 {
			end = 0
		}
		if start < 0 {
			start = 0
		}
		limitArray := (*array)[start:end]

		return reverseMessage(limitArray), total, total/limit + IsLarge(total%limit > 0)
	} else {
		return reverseMessage(*array), total, 1
	}
}

func IsLarge(res bool) int64 {
	if res {
		return 1
	} else {
		return 0
	}
}

func reverseString(input []string) []string {
	if len(input) == 0 {
		return input
	}
	return append(reverseString(input[1:]), input[0])
}

func reverseMessage(input []core.Message) []core.Message {
	if len(input) == 0 {
		return input
	}
	return append(reverseMessage(input[1:]), input[0])
}
