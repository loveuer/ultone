package sqlType

import "errors"

var (
	ErrConvertScanVal = errors.New("convert scan val to str err")
	ErrInvalidScanVal = errors.New("scan val invalid")
	ErrConvertVal     = errors.New("convert err")
)
