package core

import (
	"GoProject/BlockChain/PublicBlockChain/utils"
	"GoProject/BlockChain/PublicBlockChain/wallet"
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
)

type BlockChain struct {
	//Blocks []*Block
	// 持久化
	Client        *mongo.Client     // 数据库对象
	LastBlockHash []byte            // 保存最新区块的哈希值
}

// IsGenesisBlockExisted : 判断创世区块是否存在
func (bc *BlockChain)IsGenesisBlockExisted()bool{
	lastBlockHash:= QueryTip(bc.Client).L
	if len(lastBlockHash)!=0{
		return true
	}
	return false
}

// GetBlockChain : 获取区块链对象
func GetBlockChain(client *mongo.Client)*BlockChain{
	lastBlockHash:= QueryTip(client).L
	return &BlockChain{
		Client: client,
		LastBlockHash: lastBlockHash,
	}
}

// NewBlockChain : 创建区块链
func NewBlockChain(client *mongo.Client,address string) *BlockChain {
	var lastBlockHash []byte
	var genesisBlock *Block
	lastBlockHash= QueryTip(client).L
	if len(lastBlockHash)==0{
		// 标志区块，用于标记区块里的头指针
		tabBlock:=CreateTabBlock()
		// 加入标志区块
		InsertSensor(client,tabBlock)
		// 生成一个Coinbase交易
		txCoinbase:= NewCoinbaseTransaction(address)
		// 根据Coinbase交易生成创世区块
		genesisBlock = CreateGenesisBlock([]*Transaction{txCoinbase},tabBlock.CurBlockHash)
		InsertSensor(client,genesisBlock)
		// 设置当前最新哈希
		lastBlockHash=genesisBlock.CurBlockHash
		// 存取tip
		InsertTip(client,Tip{lastBlockHash})
		fmt.Println("Create HaHa BlockChain Successfully")
	}
	return &BlockChain{
		Client: client,
		LastBlockHash: lastBlockHash,
	}
}

//AddBlock : 添加区块到区块链中
//func (bc *BlockChain) AddBlock(data []byte) {
//	var newBlock *Block
//	height := len(bc.Blocks)
//	newBlock = NewBlock(bc.Blocks[height-1], data)
//	bc.Blocks = append(bc.Blocks, newBlock)
//}

//AddBlock : 添加区块到区块链中
func (bc *BlockChain) AddBlock(transactions []*Transaction) {
	var newBlock *Block
	tip:= QueryTip(bc.Client)
	lastBlock:= QuerySensor(bc.Client,tip.L)
	newBlock=NewBlock(lastBlock,transactions)
	// 交易签名验证
	for _,transaction:=range transactions{
		if !bc.VerifyTransaction(transaction){
			log.Panicf("Transaction verify failed!")
		}
	}
	// 保存Block
	InsertSensor(bc.Client,newBlock)
	// 更新tip
	bc.LastBlockHash=newBlock.CurBlockHash
	UpdateTip(bc.Client,bc.LastBlockHash)
}

//PrintBlockChain : 打印区块链
func (bc *BlockChain)PrintBlockChain(){
	var iterator = NewIterator(bc.Client,bc.LastBlockHash)

	for block:=iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next(){
		fmt.Println(block.String())
	}
}

// Mine : 挖矿
func (bc *BlockChain)Mine(){
	var transactions []*Transaction
	bc.AddBlock(transactions)
}

// FindUnspentUTXOs : 根据地址查找UTXO集合
/*
	从后向前遍历区块链上的每一个区块
	查找每一个交易的每一个输出
	判断交易的输出是否满足以下条件
		1、属于传入的地址
		2、是否被花费
		首先遍历区块链，查找所有的自己已花费的输出，即通过查找输入中的交易来实现，存入map中，具体通过FindSpentOutputs来实现
		然后检查每一个Vout是否已经存入该map中，如果存入则表示该输出已花费

*/
// FindUnspentUTXOs : 查找所有未花费的输出
func (bc *BlockChain)FindUnspentUTXOs(address string, transactionBuffer []*Transaction) []*UTXO {
	var unUTXOs []*UTXO
	var block *Block
	var spentTxOutputs = bc.FindSpentOutputs(address,transactionBuffer)
	// 查询缓存中的UTXO
	for _,transaction:=range transactionBuffer{
		// 遍历缓存中的交易
		Catchvout:
		for index,vout:=range transaction.Vouts{
			if vout.CheckPublicKeyWithAddress(address){
				// 验证地址
				if len(spentTxOutputs)==0{
					// 首先判断spentOutputs是否为空
					utxo:=&UTXO{
						TxHash: transaction.TransactionHash,
						Index: index,
						Output: vout,
					}
					unUTXOs=append(unUTXOs, utxo)
				}else{
					// 判断交易是否被引用
					var isUTXO bool
					for txHash,indexArray:= range spentTxOutputs{
						if txHash==hex.EncodeToString(transaction.TransactionHash) {
							isUTXO=true
							// 判断交易的输出是否被引用
							var isSpentUTXO bool
							for _,i:=range indexArray{
								// 遍历索引
								if i==index {
									//txHash==hex.EncodeToString(tx.TransactionHash) 说明当前交易已被其他交易引用
									//index==i 说明正好是当前的交易输出被其他交易引用
									// 可能append多次vout，因此要跳到指定标签，开始下一个vout
									isSpentUTXO=true
									continue Catchvout
								}
							}
							if isSpentUTXO==false{
								utxo:=&UTXO{
									TxHash: transaction.TransactionHash,
									Index: index,
									Output: vout,
								}
								unUTXOs=append(unUTXOs, utxo)
							}
						}
					}
					if isUTXO==false{
						// 如果整个交易未被引用，则整个交易的输出未被引用
						utxo:=&UTXO{
							TxHash: transaction.TransactionHash,
							Index: index,
							Output: vout,
						}
						unUTXOs=append(unUTXOs, utxo)
					}
				}
			}
		}
	}
	// 查询已在链上部分
	iterator:=NewIterator(bc.Client,bc.LastBlockHash)
	block = iterator.CurrentBlock
	for block = iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next(){
		// 遍历区块链
		for _,tx:=range block.Transactions{
			// 遍历区块中的交易
			findvout:
			for index,vout:=range tx.Vouts{
				if vout.CheckPublicKeyWithAddress(address){
					// 验证地址
					if len(spentTxOutputs)==0{
						// 首先判断spentOutputs是否为空
						utxo:=&UTXO{
							TxHash: tx.TransactionHash,
							Index: index,
							Output: vout,
						}
						unUTXOs=append(unUTXOs, utxo)
					}else{
						var isSpentOutput bool
						for txHash,indexArray:= range spentTxOutputs{
							// 遍历索引
							for _,i:=range indexArray{
								if txHash==hex.EncodeToString(tx.TransactionHash) && index==i{
									//txHash==hex.EncodeToString(tx.TransactionHash) 说明当前交易已被其他交易引用
									//index==i 说明正好是当前的输出被其他交易引用
									// 可能append多次vout，因此要跳到指定标签，开始下一个vout
									isSpentOutput=true
									continue findvout
								}
							}
						}
						if isSpentOutput==false{
							//如果未花费
							utxo:=&UTXO{
								TxHash: tx.TransactionHash,
								Index: index,
								Output: vout,
							}
							unUTXOs=append(unUTXOs, utxo)
						}
					}
				}
			}
		}
	}
	return unUTXOs
}

// FindSpentOutputs : 根据指定地址查找已花费的输出 transactionBuffer: 交易缓冲，用于缓冲多笔交易
func (bc *BlockChain) FindSpentOutputs(address string, transactionBuffer []*Transaction) map[string][]int{
	//处理address，获取publicKey
	publicKey:=utils.GetPublicKeyWithAddress(address)
	//已花费输出的缓存，用于存取索引，采用切片是存在A->b,A->c在同一个交易
	spentTxOutputs:=make(map[string][]int)
	iterator:=NewIterator(bc.Client,bc.LastBlockHash)
	var block *Block
	// 查询缓存中输入的交易
	for _,transaction:=range transactionBuffer{
		if !transaction.IsCoinbaseTransaction(){
			for _,vin:=range transaction.Vins{
				// 查询缓存中输入的交易
				if vin.CheckPublicKey(publicKey){
					key:=hex.EncodeToString(vin.TransactionHash)
					spentTxOutputs[key]=append(spentTxOutputs[key],vin.Vout)
				}
			}
		}
	}
	for block = iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next() {
		// 遍历区块链上的所有交易
		for _,tx:=range block.Transactions{
			// 遍历一个区块的所有交易
			if !tx.IsCoinbaseTransaction(){
				//排除coinbase类型交易
				for _,vin:=range tx.Vins{
					//遍历一个交易的所有输入
					if vin.CheckPublicKey(publicKey){
						key:=hex.EncodeToString(vin.TransactionHash)
						// 将输入中的十六进制transactionHash和Vout索引号作为键值添加到spentTXOutputs
						spentTxOutputs[key]=append(spentTxOutputs[key],vin.Vout)
					}
				}
			}
		}
	}
	return spentTxOutputs
}

// FindSpendableUTXOs : 查找指定地址可花费的UTXO
func (bc *BlockChain)FindSpendableUTXOs(address string,amount int,transactionBuffer []*Transaction)(int,map[string][]int){
	var value int
	//可用的UTXO
	spendableUTXOs:=make(map[string][]int)
	utxos:=bc.FindUnspentUTXOs(address,transactionBuffer)
	for _,utxo:=range utxos{
		key:=hex.EncodeToString(utxo.TxHash)
		spendableUTXOs[key]=append(spendableUTXOs[key], utxo.Index)
		value+=utxo.Output.Value
		if value>=amount {
			break
		}
	}
	if value<amount{
		fmt.Printf("Current Address %s`s balance is not enough: %d Hacoin, Current Address Balance: %d",address,amount,value)
		os.Exit(1)
	}
	return value,spendableUTXOs
}

// SignTransaction : 对交易进行签名
func (bc *BlockChain)SignTransaction(transaction *Transaction,privateKey ecdsa.PrivateKey){
	/*
		1、首先判断是否是coinbase交易
		2、遍历引用的交易，将Transaction加入
		对交易签名的基本思路：
			签名是对交易者公钥，接受者公钥以及Value的签名
			因此先对上述构成的交易副本进行哈希
			再利用椭圆曲线的Sign方法进行签名
			将上述所得签名存入Vin的Signature
	*/
	if transaction.IsCoinbaseTransaction(){
		return
	}
	quoteTransactions := make(map[string]Transaction)
	for _,vin:=range transaction.Vins{
		tx:=bc.FindTransaction(vin.TransactionHash)
		quoteTransactions[hex.EncodeToString(tx.TransactionHash)]=tx
	}
	// 签名
	transaction.Sign(privateKey,quoteTransactions)
}

// FindTransaction : 根据交易哈希查询交易
func (bc *BlockChain)FindTransaction(txHash []byte) Transaction {
	tip:= QueryTip(bc.Client)
	iterator:=NewIterator(bc.Client,tip.L)
	for block:=iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next(){
		for _,transaction:=range block.Transactions{
			if bytes.Compare(transaction.TransactionHash,txHash)==0{
				return *transaction
			}
		}
	}
	return Transaction{}
}

// VerifyTransaction : 验证交易
func (bc *BlockChain)VerifyTransaction(transaction *Transaction)bool{
	if transaction.IsCoinbaseTransaction(){
		return true
	}
	quoteTransactions:=make(map[string]Transaction)
	for _,vin:=range transaction.Vins{
		tx:=bc.FindTransaction(vin.TransactionHash)
		quoteTransactions[hex.EncodeToString(tx.TransactionHash)]=tx
	}
	return transaction.Verify(quoteTransactions)
}

// FindAllUTXO : 查找所有的UTXO
func (bc *BlockChain)FindAllUTXO()map[string][]*TxOutput {
	var utxos = make(map[string][]*TxOutput)
	spentOutputs:=bc.FindAllSpentOutput()
	iterator:=NewIterator(bc.Client,bc.LastBlockHash)
	for block:=iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next(){
		var txOutputs []*TxOutput
		for _,tx:= range block.Transactions{
			txHash:=hex.EncodeToString(tx.TransactionHash)
			WorkLoop:
			for index,vout:= range tx.Vouts{
				vins := spentOutputs[txHash]
				if len(vins)>0{
					isSpent := false
					for _,in:=range vins{
						outPublicKey:=vout.PublicKey
						inPublicKey :=in.PublicKey
						if bytes.Compare(outPublicKey, wallet.GeneratePublicKeyHash(inPublicKey))==0{
							if index == in.Vout{
								// 判断是哪个序号被引用
								isSpent=true
								continue WorkLoop
							}
						}
					}
					if isSpent==false{
						txOutputs = append(txOutputs, vout)
					}
				}else{
					txOutputs = append(txOutputs, vout)
				}

			}
			utxos[txHash]=txOutputs
		}
	}
	return utxos
}

// FindAllSpentOutput : 查找所有已花费的输出
func (bc *BlockChain)FindAllSpentOutput()map[string][]*TxInput {
	var spentOutputs = make(map[string][]*TxInput)
	iterator:=NewIterator(bc.Client,bc.LastBlockHash)
	for block:=iterator.GetCurrentBlock();iterator.HasNext();block=iterator.Next(){
		for _,tx:=range block.Transactions{
			if !tx.IsCoinbaseTransaction(){
				for _,vin:=range tx.Vins{
					txHash:=hex.EncodeToString(tx.TransactionHash)
					spentOutputs[txHash] = append(spentOutputs[txHash],vin)
				}
			}
		}
	}
	return spentOutputs
}