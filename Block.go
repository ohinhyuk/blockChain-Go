package main

type Block struct {
	PrevHash string
	Data     string
	Hash     string
}

func NewBlock(prevHash string, data string) *Block {
	block := &Block{
		PrevHash: prevHash,
		Data:     data,
	}
	block.calculateHash()
	return block
}

func (b *Block) calculateHash() {
	b.Hash = CalculateHash(b.PrevHash + b.Data)
}

func (b *Block) IsValid() bool {
	// 유효성 검증 로직 구현
	return true
}
