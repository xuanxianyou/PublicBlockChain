//交易 类型：coinbase和 转账
package core

import (
	"GoProject/BlockChain/PublicBlockChain/utils"
	wallet2 "GoProject/BlockChain/PublicBlockChain/wallet"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"log"
	"math/big"
	"strconv"
)

const(
	SystemReward="System Mine Reward"
	BonusCoin=10
)

type Transaction struct {
	TransactionHash []byte      //交易哈希（唯一标识）
	Vins            []*TxInput  //交易输入列表
	Vouts           []*TxOutput //交易输出列表
}

func (tx *Transaction)String()string{
	var transactionString string
	transactionString+="\t\t"+"TransactionHash: "+hex.EncodeToString(tx.TransactionHash)+"\n"
	// 拼接输入部分
	transactionString+="\t\tThe List of Transaction Input:\n"
	for _,txInput:=range tx.Vins{
		txHash:=hex.EncodeToString(txInput.TransactionHash)
		vout:=strconv.Itoa(txInput.Vout)
		signature:=hex.EncodeToString(txInput.Signature)
		publicKey:=hex.EncodeToString(txInput.PublicKey)
		transactionString+="\t\t\t"+"TxHash: "+txHash+"\n"
		transactionString+="\t\t\t"+"Vout: "+vout+"\n"
		transactionString+="\t\t\t"+"Signature: "+signature+"\n"
		transactionString+="\t\t\t"+"PublicKey: "+publicKey+"\n"
	}
	// 拼接输出部分
	transactionString+="\t\tThe List of Transaction Output:\n"
	for _,txOutput:=range tx.Vouts{
		value:=strconv.Itoa(txOutput.Value)
		publicKey:=hex.EncodeToString(txOutput.PublicKey)
		transactionString+="\t\t\t"+"Value: "+value+"\n"
		transactionString+="\t\t\t"+"PublicKey: "+publicKey+"\n"
	}
	return transactionString
}
//NewCoinbaseTransaction : 挖矿交易
func NewCoinbaseTransaction(address string) *Transaction {
	publicKey:=utils.GetPublicKeyWithAddress(address)
	var transaction *Transaction
	// 交易输入
	txInput:=&TxInput{
		TransactionHash: []byte{},
		Vout: -1,
		Signature: nil,
		PublicKey: nil,
	}
	// 交易输出
	txOutput:=&TxOutput{
		Value:     BonusCoin,
		PublicKey: publicKey,
	}
	// 交易
	transaction=&Transaction{
		TransactionHash: nil,
		Vins: []*TxInput{txInput},
		Vouts: []*TxOutput{txOutput},
	}
	//设置哈希
	transaction.SetTransactionHash()
	return transaction
}

// SetTransactionHash : 设置交易的Hash
func (tx *Transaction)SetTransactionHash(){
	var result bytes.Buffer
	encoder:=gob.NewEncoder(&result)
	err:=encoder.Encode(tx)
	if err!=nil{
		log.Println(err)
	}
	//生成哈希值
	hash:=sha256.Sum256(result.Bytes())
	tx.TransactionHash=hash[:]
}

// NewTransferTransaction : 实现转账交易
func NewTransferTransaction(from string,to string,amount int,blockchain *BlockChain,transactionBuffer []*Transaction)*Transaction {
	// 获取钱包
	wallets:= wallet2.NewWallets()
	wallet:=wallets.WalletsMap[from]
	// 获取查找到的money和UTXOs集合
	money,spendableUTXOs:=blockchain.FindSpendableUTXOs(from,amount,transactionBuffer)
	// 构造输入交易集合
	// 输入
	var txInputs []*TxInput
	for transactionHashString,indexArray:=range spendableUTXOs{
		transactionHash,err:=hex.DecodeString(transactionHashString)
		if err!=nil{
			log.Panicf("transactionHash string decode to hash ERROR:%v\n", err)
		}
		for _,index:=range indexArray{
			txInput:= NewTxInput(transactionHash,index,from)
			txInputs=append(txInputs, txInput)
		}
	}
	// 输出
	txOutput:= NewTxOutput(amount,to)
	var txOutputs []*TxOutput
	txOutputs=append(txOutputs, txOutput)
	// 找零
	 if amount < money{
	 	txOutput:= NewTxOutput(money-amount,from)
		 txOutputs=append(txOutputs, txOutput)
	 }
	transaction:=&Transaction{
		TransactionHash: nil,
		Vins: txInputs,
		Vouts: txOutputs,
	}
	transaction.SetTransactionHash()
	blockchain.SignTransaction(transaction,wallet.PrivateKey)
	return transaction
}

// IsCoinbaseTransaction : 判断是否是Coinbase类型交易
func (tx *Transaction)IsCoinbaseTransaction() bool {
	return tx.Vins[0].Vout==-1 && len(tx.Vins[0].TransactionHash)==0
}

// Sign : 对交易进行签名
func (tx *Transaction)Sign(privateKey ecdsa.PrivateKey,quoteTransactions map[string]Transaction){
	for _,vin:=range tx.Vins{
		if quoteTransactions[hex.EncodeToString(vin.TransactionHash)].TransactionHash==nil{
			// 判断交易是否在引用的交易中，如果都在，则说明未篡改
			log.Panicf("Sign transaction failed!\n")
		}
	}
	txCopy:=tx.TransactionCopy()
	for vinId,vin:=range txCopy.Vins{
		// 获取被引用的交易
		quoteTx:=quoteTransactions[hex.EncodeToString(vin.TransactionHash)]
		// 找到发送者
		txCopy.Vins[vinId].PublicKey=quoteTx.Vouts[vin.Vout].PublicKey
		//
		hash:=sha256.Sum256(txCopy.Serialize())
		txCopy.TransactionHash=hash[:]
		// 调用签名函数
		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txCopy.TransactionHash)
		if err!=nil{
			log.Panicf("Transaction sign Error:%v\n",err)
		}
		signature:=append(r.Bytes(),s.Bytes()...)
		tx.Vins[vinId].Signature=signature
	}
}

// Verify : 对交易签名进行验证
func (tx *Transaction)Verify(quoteTransactions map[string]Transaction)bool{
	for _,vin:=range tx.Vins{
		if quoteTransactions[hex.EncodeToString(vin.TransactionHash)].TransactionHash == nil{
			// 检查交易是否在引用的交易中，如果都在，则说明未篡改
			log.Panicf("Transaction verify failed!\n")
		}
	}
	txCopy:=tx.TransactionCopy()
	curve:=elliptic.P256()
	for vinId,vin:=range tx.Vins{
		quoteTx:=quoteTransactions[hex.EncodeToString(vin.TransactionHash)]
		txCopy.Vins[vinId].PublicKey=quoteTx.Vouts[vin.Vout].PublicKey
		//hash
		hash:= sha256.Sum256(txCopy.Serialize())
		txCopy.TransactionHash = hash[:]
		// 求publicKey
		x:=&big.Int{}
		y:=&big.Int{}
		x.SetBytes(vin.PublicKey[:len(vin.PublicKey)/2])
		y.SetBytes(vin.PublicKey[len(vin.PublicKey)/2:])
		publicKey:=&ecdsa.PublicKey{
			Curve:curve,
			X: x,
			Y: y,
		}
		// 求r,s
		r:=&big.Int{}
		s:=&big.Int{}
		r.SetBytes(vin.Signature[:len(vin.Signature)/2])
		s.SetBytes(vin.Signature[len(vin.Signature)/2:])
		// 验证
		if !ecdsa.Verify(publicKey,txCopy.TransactionHash,r,s){
			return false
		}

	}
	return true

}

// TransactionCopy : 生成交易副本
func (tx *Transaction)TransactionCopy()*Transaction {
	/*
		生成交易的副本，交易副本只包含原交易的部分属性
		采用副本主要考虑到:
			1、签名时只用到以上部分属性
			2、采用部分属性可以生成新的TransactionHash
	*/
	var inputs []*TxInput
	var outputs []*TxOutput
	// 输入
	for _,vin:=range tx.Vins{
		inputs=append(inputs, &TxInput{vin.TransactionHash,vin.Vout,nil,nil})
	}
	// 输出
	for _,vout:=range tx.Vouts{
		outputs= append(outputs, &TxOutput{vout.Value,vout.PublicKey})
	}
	transaction:=&Transaction{tx.TransactionHash,inputs,outputs}
	return transaction
}

//
func (tx *Transaction)Serialize()[]byte{
	var transactionBytes bytes.Buffer
	encoder:=gob.NewEncoder(&transactionBytes)
	if err:=encoder.Encode(tx);err!=nil{
		log.Panicf("Transaction serialize Error:%v\n",err)
	}
	return transactionBytes.Bytes()
}
