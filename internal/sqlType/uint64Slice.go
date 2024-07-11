package sqlType

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cast"
)

type NumSlice[T ~int | ~int64 | ~uint | ~uint64] []T

func (n *NumSlice[T]) Scan(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return ErrConvertScanVal
	}

	length := len(str)

	if length <= 0 {
		*n = make(NumSlice[T], 0)
		return nil
	}

	if str[0] != '{' || str[length-1] != '}' {
		return ErrInvalidScanVal
	}

	str = str[1 : length-1]
	if len(str) == 0 {
		*n = make(NumSlice[T], 0)
		return nil
	}

	numStrs := strings.Split(str, ",")
	nums := make([]T, len(numStrs))

	for idx := range numStrs {
		num, err := cast.ToInt64E(strings.TrimSpace(numStrs[idx]))
		if err != nil {
			return fmt.Errorf("%w: can't convert to %T", ErrConvertVal, T(0))
		}

		nums[idx] = T(num)
	}

	*n = nums

	return nil
}

func (n NumSlice[T]) Value() (driver.Value, error) {
	if n == nil {
		return "{}", nil
	}

	if len(n) == 0 {
		return "{}", nil
	}

	ss := make([]string, 0, len(n))
	for idx := range n {
		ss = append(ss, strconv.Itoa(int(n[idx])))
	}

	s := strings.Join(ss, ", ")

	return fmt.Sprintf("{%s}", s), nil
}
