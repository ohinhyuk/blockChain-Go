package main

type Blockchain struct {
	Chain []Block
}

func NewBlockchain() *Blockchain {
	genesisBlock := NewBlock("0", "Genesis Block")
	return &Blockchain{
		Chain: []Block{*genesisBlock},
	}
}

func (bc *Blockchain) AddBlock(newBlock Block) {
	if newBlock.IsValid() && bc.IsChainValid() {
		bc.Chain = append(bc.Chain, newBlock)
	}
}

func (bc *Blockchain) IsChainValid() bool {
	for i := 1; i < len(bc.Chain); i++ {
		currentBlock := bc.Chain[i]
		prevBlock := bc.Chain[i-1]

		if currentBlock.Hash != CalculateHash(currentBlock.PrevHash+currentBlock.Data) {
			return false
		}

		if currentBlock.PrevHash != prevBlock.Hash {
			return false
		}
	}
	return true
}
