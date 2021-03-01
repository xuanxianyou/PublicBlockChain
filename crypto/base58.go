package crypto

import (
	//"GoProject/BlockChain/PublicBlockChain/core"
	"bytes"
	"math/big"
)

// base58 编码基数表
var base58Alphabet = []byte("123456789" +
	"abcdefghijkmnopqrstuvwxyz" +
	"ABCDEFGHJKLMNPQRSTUVWXYZ")


func Base58Encode(input []byte) []byte{
	// 将字节数组转化为big.Int，不断对58取余，将获得的余数索引从base58Alphabet获得对应的编码
	var result []byte
	x:=big.NewInt(0).SetBytes(input)
	base:=big.NewInt(int64(len(base58Alphabet)))
	zero:=big.NewInt(0)
	mod:=&big.Int{}
	for x.Cmp(zero) !=0 {
		x,mod=x.DivMod(x,base,mod)
		result=append(result,base58Alphabet[mod.Int64()])
	}
	// 反转result
	Reserve(result)
	// 添加一个前缀 1，标识这是一个地址
	result=append([]byte{'1'}, result...)
	return result
}

func Reserve(slice []byte){
	for i,j:=0,len(slice)-1;i < j;i,j=i+1,j-1{
		slice[i],slice[j]=slice[j],slice[i]
	}
}

func Base58Decode(address []byte) []byte{
	var result =big.NewInt(0)
	// 去掉版本前缀
	data:=address[1:]
	for _,b:=range data{
		byteIndex:=bytes.IndexByte(base58Alphabet,b)
		result.Mul(result,big.NewInt(58))
		result.Add(result,big.NewInt(int64(byteIndex)))
	}
	decoded := result.Bytes()
	return decoded
}