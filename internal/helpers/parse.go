package helpers

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func SafeInt(value string, defaultVal int64) int64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return defaultVal
	}
	clean := strings.ReplaceAll(value, ",", "")
	if n, err := strconv.ParseInt(clean, 10, 64); err == nil {
		return n
	}
	return defaultVal
}

func ParseAmountArg(arg string, balance int64) (int, error) {
	switch strings.ToLower(arg) {
	case "all", "a":
		return int(balance), nil
	case "half", "h":
		return int(balance / 2), nil
	default:
		amount := ParseAmount(arg)
		if amount <= 0 {
			return 0, fmt.Errorf("invalid amount")
		}
		return int(amount), nil
	}
}

func ParseAmount(amount string) int64 {
	amount = strings.ToLower(strings.TrimSpace(amount))
	if amount == "" {
		return 0
	}

	var mult int64
	switch amount[len(amount)-1] {
	case 'k':
		mult = 1_000
	case 'm':
		mult = 1_000_000
	case 'b':
		mult = 1_000_000_000
	case 't':
		mult = 1_000_000_000_000
	case 'q':
		mult = 1_000_000_000_000_000
	case 'z':
		return math.MaxInt64
	default:
		return SafeInt(amount, 0)
	}

	base := SafeInt(amount[:len(amount)-1], 0)
	return safeMulClamp(base, mult)
}

func safeMulClamp(a, b int64) int64 {
	if a == 0 || b == 0 {
		return 0
	}
	neg := (a < 0) != (b < 0)
	ua := abs(a)
	ub := abs(b)
	if ua > math.MaxInt64/ub {
		if neg {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	res := ua * ub
	if neg {
		return -res
	}
	return res
}

func abs(x int64) int64 {
	if x < 0 {
		if x == math.MinInt64 {
			return math.MaxInt64
		}
		return -x
	}
	return x
}
