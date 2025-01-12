package translit

// CharCategory represents character types in transliteration schemes
type CharCategory int

const (
	CategoryNone CharCategory = iota
	CategoryVowel
	CategoryConsonant
	CategoryDigit
	CategoryOther
)

// Context tracks transliteration state
type Context struct {
	LastChar      rune
	Category      CharCategory
	InEnglishMode bool
	HasVirama     bool
}

func NewContext() *Context {
	return &Context{
		Category: CategoryNone,
	}
}

func (c *Context) ToggleEnglishMode() {
	c.InEnglishMode = !c.InEnglishMode
}

func (c *Context) UpdateWithOutput(output string) {
	if len(output) > 0 {
		c.LastChar = []rune(output)[0]
		// TODO: Update category based on character properties
	}
}
