package utils

type PascalShortString struct {
	capacity int
	length   int
	bytes    []byte
}

func NewPascalShortString(content string, capacity int) *PascalShortString {
	// TODO: Should all bytes on the slice be 0 to start with?
	bytes := make([]byte, capacity)
	contentBytes := []byte(content)
	for i := 0; i < len(content); i++ {
		if i == capacity {
			break
		}
		bytes[i] = contentBytes[i]
	}
	return &PascalShortString{
		capacity: capacity,
		length:   len(content),
		bytes:    bytes,
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
	contentBytes := []byte(content)
	if len(content) > 0 {
		for i := 0; i < len(content); i++ {
			if i == pss.capacity {
				break
			}
			pss.bytes[i] = contentBytes[i]
		}
	}
	if len(content) <= pss.capacity {
		pss.length = len(content)
	} else {
		pss.length = pss.capacity
	}
}
