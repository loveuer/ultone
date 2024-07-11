package sqlType

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
)

type StrSlice []string

func (s *StrSlice) Scan(val interface{}) error {

	str, ok := val.(string)
	if !ok {
		return ErrConvertScanVal
	}

	if len(str) < 2 {
		return nil
	}

	bs := make([]byte, 0, 128)
	bss := make([]byte, 0, 2*len(str))

	quoteCount := 0

	for idx := 1; idx < len(str)-1; idx++ {
		// 44: ,     92: \      34: "
		quote := str[idx]
		switch quote {
		case 44:
			if quote == 44 && str[idx-1] != 92 && quoteCount == 0 {
				if len(bs) > 0 {
					if !(bs[0] == 34 && bs[len(bs)-1] == 34) {
						bs = append([]byte{34}, bs...)
						bs = append(bs, 34)
					}

					bss = append(bss, bs...)
					bss = append(bss, 44)
				}
				bs = bs[:0]
			} else {
				bs = append(bs, quote)
			}
		case 34:
			if str[idx-1] != 92 {
				quoteCount = (quoteCount + 1) % 2
			}
			bs = append(bs, quote)
		default:
			bs = append(bs, quote)
		}

		//bs = append(bs, str[idx])
	}

	if len(bs) > 0 {
		if !(bs[0] == 34 && bs[len(bs)-1] == 34) {
			bs = append([]byte{34}, bs...)
			bs = append(bs, 34)
		}

		bss = append(bss, bs...)
	} else {
		if len(bss) > 2 {
			bss = bss[:len(bss)-2]
		}
	}

	bss = append([]byte{'['}, append(bss, ']')...)

	if err := json.Unmarshal(bss, s); err != nil {
		return err
	}

	return nil
}

func (s StrSlice) Value() (driver.Value, error) {
	if s == nil {
		return "{}", nil
	}

	if len(s) == 0 {
		return "{}", nil
	}

	buf := &bytes.Buffer{}

	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(s); err != nil {
		return "{}", err
	}

	bs := buf.Bytes()

	bs[0] = '{'

	if bs[len(bs)-1] == 10 {
		bs = bs[:len(bs)-1]
	}

	bs[len(bs)-1] = '}'

	return string(bs), nil
}
