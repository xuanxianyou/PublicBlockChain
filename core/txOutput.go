package core

import (
	"GoProject/BlockChain/PublicBlockChain/utils"
	"bytes"
)

type TxOutput struct {
	Value      int      //金额
	PublicKey  []byte   //UTXO的拥有者，在这里的PublicKey并不是公钥，而是PublicKeyHash
}

func NewTxOutput(value int,address string)*TxOutput {
	publicKey:=utils.GetPublicKeyWithAddress(address)
	return &TxOutput{
		Value: value,
		PublicKey: publicKey,
	}
}
//func (txOutput *TxOutput)CheckScriptPublicKeyWithAddress(address string)bool{
//	return address==txOutput.ScriptPublicKey
//}

// CheckPublicKeyWithAddress : 根据地址检查交易输出的PublicKey
func (txOutput *TxOutput)CheckPublicKeyWithAddress(address string)bool{
	publicKey:=utils.GetPublicKeyWithAddress(address)
	return bytes.Compare(publicKey,txOutput.PublicKey)==0
}
