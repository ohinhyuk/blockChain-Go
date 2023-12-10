package main

// Blockchain keeps a sequence of Blocks
// type Blockchain struct {
// 	blocks []*Block
// }
type Blockchain struct {
	Blocks []*Block
	Nodes  []*Node
}

// AddBlock saves provided data as a block in the blockchain
func (bc *Blockchain) AddBlock(data string) *Block {
	prevBlock := bc.Blocks[len(bc.Blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.Blocks = append(bc.Blocks, newBlock)

	return newBlock
}

// NewBlockchain creates a new Blockchain with genesis Block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}, []*Node{}}
}