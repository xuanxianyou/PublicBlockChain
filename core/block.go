package core

import (
	"GoProject/BlockChain/PublicBlockChain/utils"
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"strconv"
	"time"
)

// Block : 区块
type Block struct {
	TimeStamp    int64  //区块时间戳
	CurBlockHash []byte //当前区块哈希
	PreBlockHash []byte //前区块哈希
	Height       int64  //区块索引、区块号
	//Data       []byte //交易数据
	MerkleRoot   []byte         //Merkle树根节点哈希
	Transactions []*Transaction //交易
	Nonce        int64          //pow运行时动态修改的数据
}

//NewBlock : 新建区块
func NewBlock(preBlock *Block, transactions []*Transaction) *Block {
	var block *Block
	block = &Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: preBlock.CurBlockHash,
		Height:       preBlock.Height + 1,
		Transactions: transactions,
	}
	// block.SetCurBlock()
	// 工作量证明
	pow := NewProofOfWork(block)
	hash, nonce := pow.Run()
	//设置CurBlockHash和Nonce
	block.CurBlockHash = hash[:]
	block.Nonce = int64(nonce)
	return block
}

//SetCurBlock : 设置区块哈希
func (block *Block) SetCurBlockHash() {
	timeStampBytes := utils.IntToHex(block.TimeStamp)
	heightBytes := utils.IntToHex(block.Height)
	// 通过byte二维数组转化成一维byte数组
	byteBlocks := bytes.Join([][]byte{
		timeStampBytes,
		block.PreBlockHash,
		heightBytes,
		block.TransactionsSerialize(),
	}, []byte{})
	// 对一维数组进行sha256运算
	hash := sha256.Sum256(byteBlocks)

	block.CurBlockHash = hash[:]
}

// String : block 转 string
func (block *Block) String() string {
	t := time.Unix(block.TimeStamp, 0).String()
	curBlockHash := hex.EncodeToString(block.CurBlockHash)
	preBlockHash := hex.EncodeToString(block.PreBlockHash)
	height := strconv.FormatInt(block.Height, 10)
	nonce := strconv.FormatInt(block.Nonce, 10)
	var transactionsString string
	for _,transaction:=range block.Transactions{
		transactionsString+=transaction.String()
	}
	return "\theight: " + height + "\n" +
		"\tnonce: " + nonce + "\n" +
		"\tcreateTime: " + t + "\n" +
		"\tpreBlockHash: " + preBlockHash + "\n" +
		"\tcurBlockHash: " + curBlockHash + "\n" +
	    "\ttransactions: \n" + transactionsString
}

//Serialize : Block 序列化
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	//新建编码对象
	encoder := gob.NewEncoder(&buffer)
	//编码序列化
	if err := encoder.Encode(block); err != nil {
		log.Panic(err.Error())
	}
	return buffer.Bytes()
}

//反序列化
func Deserialize(blockBytes []byte) *Block {
	var block Block
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	if err := decoder.Decode(&block); err != nil {
		log.Panic(err.Error())
	}
	return &block
}

// CreateGenesisBlock : 生成创世区块
func CreateGenesisBlock(transactions []*Transaction,tabBlockHash []byte) *Block {
	block := Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: tabBlockHash,
		Height:       0,
		Transactions: transactions,
	}
	//block.SetCurBlockHash()
	// 工作量证明
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()
	//设置CurBlockHash和Nonce
	block.CurBlockHash = hash[:]
	block.Nonce = int64(nonce)
	return &block
}

// CreateTabBlock : 创建标志区块
func CreateTabBlock()*Block{
	block := &Block{
		TimeStamp:    time.Now().Unix(),
		PreBlockHash: []byte(""),
		Height:       -1,
		Transactions: []*Transaction{},
	}
	block.SetCurBlockHash()
	return block
}

// TransactionsSerialize : 将区块中的交易集合序列化
func (block *Block) TransactionsSerialize() []byte {
	var transactionsBytes [][]byte
	for _, transaction := range block.Transactions {
		transactionsBytes = append(transactionsBytes, transaction.TransactionHash)
	}
	if len(transactionsBytes)==0{
		transactionHash := sha256.Sum256(bytes.Join(transactionsBytes, []byte{}))
		return transactionHash[:]
	}
	//fmt.Println(transactionsBytes)
	merkleTree:=NewMerkleTree(transactionsBytes)
	return merkleTree.Root.Data
}
