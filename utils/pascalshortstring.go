package utils

type PascalShortString struct {
	capacity int
	length   int
	bytes    []byte
}

func NewPascalShortString(content string, capacity int) *PascalShortString {
	return &PascalShortString{
		capacity: capacity,
		length:   len(content),
		bytes:    []byte(content),
	}
}

func (pss *PascalShortString) String() string {
	var result = pss.bytes[:pss.length]
	return string(result)
}

func (pss *PascalShortString) Raw() []byte {
	return pss.bytes
}

func (pss *PascalShortString) Len() int {
	return pss.length
}

func (pss *PascalShortString) Capacity() int {
	return pss.capacity
}

func (pss *PascalShortString) NewContent(content string) {
	var bytes = []byte(content)
	for i := 0; i < pss.capacity; i++ {
		pss.bytes[i] = bytes[i]
	}
}
