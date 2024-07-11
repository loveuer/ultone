package interfaces

type Enum interface {
	Value() int64
	Code() string
	Label() string

	MarshalJSON() ([]byte, error)

	All() []Enum
}
