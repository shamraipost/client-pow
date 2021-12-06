package proofofwork

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"math"
	"math/big"
	"strconv"
)

const maxNonce = math.MaxInt64

type ProofOfWork struct {
	HashCash   string
	TargetBits int
}

func IntToHex(n int64) []byte {
	return []byte(strconv.FormatInt(n, 16))
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			[]byte(pow.HashCash),
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Search() (int, []byte) {
	target := big.NewInt(1)
	target.Lsh(target, uint(160-pow.TargetBits))
	var sha1Hash []byte
	var hashInt big.Int
	nonce := 0
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash := sha1.New()
		hash.Write(data)
		sha1Hash = hash.Sum(nil)
		hashInt.SetBytes(sha1Hash[:])
		if hashInt.Cmp(target) == -1 {
			break
		}
		nonce++
	}
	fmt.Printf("iteration: %d, found: %x \n", nonce, sha1Hash)
	return nonce, sha1Hash
}
