package model

type Enum interface {
	Value() int64
	Code() string
	Label() string

	MarshalJSON() ([]byte, error)

	All() []Enum
}

type OpLogger interface {
	Enum
	Render(content map[string]any) (string, error)
	Template() string
}
