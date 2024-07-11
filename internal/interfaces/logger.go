package interfaces

type OpLogger interface {
	Enum
	Render(content map[string]any) (string, error)
	Template() string
}
