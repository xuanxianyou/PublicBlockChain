package core

import (
	"encoding/hex"
)

// UTXOSet : 区块链上的utxo集合
type UTXOSet struct{
	Blockchain *BlockChain
}

// ResetUTXOSet : 重置
func (utxoSet *UTXOSet)ResetUTXOSet(){
	ResetUTXOTable(utxoSet.Blockchain)
}

// GetBalance : 查询余额
func (utxoSet *UTXOSet)GetBalance(address string) int {
	var amount int
	utxos:= FindUTXOTable(utxoSet.Blockchain,address)
	for _,utxo:=range utxos{
		amount+=utxo.Output.Value
	}
	return amount
}

// UpdateUTXOTable : 更新UTXO Table，以使数据库UTXO正确，正确获取Balance
func (utxoSet *UTXOSet)UpdateUTXOTable(){
	iterator:= NewIterator(utxoSet.Blockchain.Client,utxoSet.Blockchain.LastBlockHash)
	lastBlock:=iterator.GetCurrentBlock()
	for _,tx:=range lastBlock.Transactions{
		// 因为每一次区块上链，UTXO Table都更新一次，因此只需要只查找最近一个区块的交易
		if !tx.IsCoinbaseTransaction(){
			for _,vin:=range tx.Vins{
				// 将已经被当前交易引用的UTXO删掉
				updataOutputs:=&TxOutputs{}
				outs:= QueryTxOutputs(utxoSet.Blockchain.Client,hex.EncodeToString(vin.TransactionHash))
				for outputIndex,out:=range outs.TxOutputs{
					if vin.Vout!=outputIndex{
						updataOutputs.TxOutputs=append(updataOutputs.TxOutputs, out)
					}
				}
				if len(updataOutputs.TxOutputs)==0{
					// 如果交易中没有UTXO，删除该交易
					DeleteUTXOTable(utxoSet.Blockchain.Client,hex.EncodeToString(vin.TransactionHash))
				}else{
					// 将更新后的UTXO数据存入数据库
					UpdateUTXOTable(utxoSet.Blockchain.Client,hex.EncodeToString(vin.TransactionHash),updataOutputs.TxOutputs)
				}
			}
		}
		newOutputs:= TxOutputs{}
		newOutputs.TxOutputs=append(newOutputs.TxOutputs, tx.Vouts...)
		InsertUTXOTable(utxoSet.Blockchain.Client,hex.EncodeToString(tx.TransactionHash),newOutputs.TxOutputs)
	}
}