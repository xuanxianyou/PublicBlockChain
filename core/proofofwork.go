package core

import (
	"GoProject/BlockChain/PublicBlockChain/utils"
	"bytes"
	"crypto/sha256"
	"math/big"
)

type ProofOfWork struct {
	//需要共识验证的区块
	Block *Block
	//目标难度哈希
	Target *big.Int

}


// NewProofOfWork : 创建ProofOfWork Worker
func NewProofOfWork(block *Block) *ProofOfWork {
	// 0000248870362186db60982d1ce3a7e4e6c05226516a872d779988437aea44c5
	//初始化nonce
	target:=big.NewInt(1)
	//nonce左移运算
	target.Lsh(target,256-TargetBit)
	return &ProofOfWork{
		Block: block,
		Target: target,
	}
}

//Run : 执行ProofOfWork
func (pow *ProofOfWork)Run() ([32]byte, int){
	var hash [32]byte
	var nonce = 0
	var hashInt big.Int
	//计算符合条件的Hash
	for{
		blockBytes:=pow.prepareData(int64(nonce))
		hash = sha256.Sum256(blockBytes)
		hashInt.SetBytes(hash[:])
		// 检测生成的哈希值是否符合条件：pow.Target > hashInt
		if pow.Target.Cmp(&hashInt)==1{
			//找到了目标哈希值
			break
		}else{
			nonce++
		}

	}
	return hash,nonce
}

//prepareData : 准备ProofOfWork运算的数据
func (pow *ProofOfWork)prepareData(nonce int64)[]byte{
	timeStampBytes := utils.IntToHex(pow.Block.TimeStamp)
	heightBytes := utils.IntToHex(pow.Block.Height)

	// 通过byte二维数组转化成一维byte数组
	blockBytes := bytes.Join([][]byte{
		timeStampBytes,
		pow.Block.PreBlockHash,
		heightBytes,
		pow.Block.TransactionsSerialize(),
		utils.IntToHex(nonce),
		utils.IntToHex(TargetBit),
	}, []byte{})
	return blockBytes
}