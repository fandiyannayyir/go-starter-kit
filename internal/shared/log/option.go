package log

type Option struct {
	OutputPath string
	Level      string
}

func NewDefaultOption() *Option {
	return &Option{
		OutputPath: "application.log",
		Level:      "info",
	}
}
