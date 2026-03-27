package domain

type OutputFmt uint

const (
	DefaultFmt OutputFmt = iota
	TableFmt
	HTMLFmt
)

type RunConfig struct {
	CurrentRepo  *string
	Fmt          OutputFmt
	HTMLTitle    string
	HTMLTemplate string
}
