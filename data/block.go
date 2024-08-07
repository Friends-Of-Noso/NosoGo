package data

type Block struct {
	Number uint64
}

func NewBlockZero() *Block {
	return &Block{
		Number: 0,
	}
}
