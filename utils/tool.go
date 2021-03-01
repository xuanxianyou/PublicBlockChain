package utils

import (
	"GoProject/BlockChain/PublicBlockChain/crypto"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"log"
)

func IntToHex(data int64)[]byte{
	buffer:=new(bytes.Buffer)
	err:=binary.Write(buffer,binary.BigEndian,data)
	if err!=nil{
		log.Panicf("int transact to []byte failed! %v\n",err)
	}
	return buffer.Bytes()
}


func JsonToSlice(jsonString string)[]string{
	var slice []string
	if err:=json.Unmarshal([]byte(jsonString),&slice);err!=nil{
		log.Panic(err)
	}
	return slice
}


// GetPublicKeyWithAddress : 根据地址获取交易输出的PublicKey，从这里应该知道为什么用PublicKeyHash代替PublicKey
func GetPublicKeyWithAddress(address string)[]byte{
	addressBytes:=[]byte(address)
	decodedHash:=crypto.Base58Decode(addressBytes)
	publicKeyHash:=decodedHash[:len(decodedHash)-32]
	return publicKeyHash
}