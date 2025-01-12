package translit

// Result represents the outcome of a transliteration operation
type Result struct {
	Output         string
	BackspaceCount int
	Error          error
}
