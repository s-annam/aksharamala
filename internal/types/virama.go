package types

type ViramaMode int

const (
	UnknownMode ViramaMode = iota
	NormalMode
	SmartMode
)

var viramaModeStrings = map[string]ViramaMode{
	"normal": NormalMode,
	"smart":  SmartMode,
}

func ParseViramaMode(mode string) ViramaMode {
	if v, ok := viramaModeStrings[mode]; ok {
		return v
	}
	return UnknownMode
}

func (v ViramaMode) String() string {
	for k, value := range viramaModeStrings {
		if value == v {
			return k
		}
	}
	return "unknown"
}
