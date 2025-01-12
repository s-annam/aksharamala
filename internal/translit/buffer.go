package translit

// Buffer manages input character buffering
type Buffer struct {
	content []rune
}

func NewBuffer() *Buffer {
	return &Buffer{
		content: make([]rune, 0, 32),
	}
}

func (b *Buffer) Append(r rune)  { b.content = append(b.content, r) }
func (b *Buffer) Clear()         { b.content = b.content[:0] }
func (b *Buffer) Len() int       { return len(b.content) }
func (b *Buffer) String() string { return string(b.content) }
func (b *Buffer) First() rune    { return b.content[0] }
func (b *Buffer) RemoveFirst()   { b.content = b.content[1:] }
func (b *Buffer) Remove(n int)   { b.content = b.content[n:] }
