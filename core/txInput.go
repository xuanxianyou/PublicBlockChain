package core

import (
	"GoProject/BlockChain/PublicBlockChain/wallet"
	"bytes"
)

type TxInput struct{
	TransactionHash []byte    // 交易哈希，不指当前交易的哈希
	Vout            int       // 引用上一笔交易输出的索引号
	Signature       []byte    // 数字签名
	PublicKey       []byte    // 发送者的公钥,用于脚本签名
}

func NewTxInput(transactionHash []byte,vout int,address string)*TxInput {
	wallets:= wallet.NewWallets()
	wallet:=wallets.WalletsMap[address]
	publicKey:=wallet.PublicKey
	txInput:=&TxInput{
		TransactionHash: transactionHash,
		Vout: vout,
		PublicKey: publicKey,
	}
	return txInput
}
//func (txInput *TxInput)CheckScriptSignatureWithAddress(address string)bool{
//	return address==txInput.ScriptSignature
//}

func (txInput *TxInput)CheckPublicKey(publicKeyHash []byte)bool{
	/*
		因为在TxOutput中的PublicKey是经过sha256和ripemd160运算后的PublicKeyHash
		所以需要对txInputPublicKey进行上述运算后才能比较
	 */
	hash:= wallet.GeneratePublicKeyHash(txInput.PublicKey)

	return bytes.Compare(hash,publicKeyHash)==0
}