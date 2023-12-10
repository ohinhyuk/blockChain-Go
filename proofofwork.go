package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

var (
	maxNonce = math.MaxInt64
)

const targetBits = 16

// ProofOfWork represents a proof-of-work
type ProofOfWork struct {
	block  *Block
	target *big.Int
}

// NewProofOfWork builds and returns a ProofOfWork
func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.block.PrevBlockHash,
			pow.block.Data,
			IntToHex(pow.block.Timestamp),
			IntToHex(int64(targetBits)),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

// Run performs a proof-of-work
func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0


	for nonce < maxNonce { 					// nonce 값이 maxNonce 보다 작을 때까지 반복
		data := pow.prepareData(nonce)

		hash = sha256.Sum256(data) 			// 해시값 계산
		fmt.Printf("\r%x", hash) 
		hashInt.SetBytes(hash[:]) 

		if hashInt.Cmp(pow.target) == -1 { 	// 조건을 만족하는 해시값 찾음
			fmt.Print("\n\n")
			return nonce, hash[:] 
		} else {
			nonce++ 						// nonce 값 증가
		}
	}
	fmt.Print("\n\n")
	return -1, nil // 채굴 실패
}

// Validate validates block's PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}