package filtering

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

func coerceValue[T any](value any) (any, error) {
	var zero T

	switch any(zero).(type) {
	case time.Time:
		return coerceTime(value)
	default:
		return coerceScalar[T](value)
	}
}

func coerceTime(value any) (any, error) {
	var err error

	switch v := value.(type) {
	case Any[string]:
		var ts []time.Time
		if ts, err = mapErr(v, parseTimestamp); err == nil {
			value = pq.Array(ts)
		}
	case string:
		value, err = parseTimestamp(v)
	default:
		err = errTimestamp
	}

	return value, err
}

func coerceScalar[T any](value any) (any, error) {
	var zero T

	switch value.(type) {
	case Any[T]:
		return pq.Array(value), nil
	case T:
		return value, nil
	default:
		switch any(zero).(type) {
		case int:
			return coerceNumber[int](value)
		case int8:
			return coerceNumber[int8](value)
		case int16:
			return coerceNumber[int16](value)
		case int32:
			return coerceNumber[int32](value)
		case int64:
			return coerceNumber[int64](value)
		case uint:
			return coerceNumber[uint](value)
		case uint8:
			return coerceNumber[uint8](value)
		case uint16:
			return coerceNumber[uint16](value)
		case uint32:
			return coerceNumber[uint32](value)
		case uint64:
			return coerceNumber[uint64](value)
		case float32:
			return coerceNumber[float32](value)
		case float64:
			return coerceNumber[float64](value)
		default:
			return nil, fmt.Errorf("value must be %T", zero)
		}
	}
}

type number interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

func coerceNumber[T number](value any) (T, error) {
	switch v := value.(type) {
	case int:
		return T(v), nil
	case int8:
		return T(v), nil
	case int16:
		return T(v), nil
	case int32:
		return T(v), nil
	case int64:
		return T(v), nil
	case uint:
		return T(v), nil
	case uint8:
		return T(v), nil
	case uint16:
		return T(v), nil
	case uint32:
		return T(v), nil
	case uint64:
		return T(v), nil
	case float32:
		return T(v), nil
	case float64:
		return T(v), nil
	default:
		return 0, fmt.Errorf("value must be %T", T(0))
	}
}

func parseTimestamp(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, errTimestamp
	}

	return t, err
}

var errTimestamp = fmt.Errorf("value must be RFC3339 timestamp")

func mapErr[T any, U any](
	items []T,
	iteratee func(T) (U, error),
) (out []U, err error) {
	out = make([]U, len(items))

	for i, item := range items {
		out[i], err = iteratee(item)
		if err != nil {
			return nil, err
		}
	}

	return out, nil
}
