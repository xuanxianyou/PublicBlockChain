package wallet

import (
	"GoProject/BlockChain/PublicBlockChain/crypto"
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"golang.org/x/crypto/ripemd160"
	"log"
)

/*
	钱包是一个公私钥对
	钱包的基本结构：
		私钥    esdsa.PrivateKey: 椭圆曲线加密
		公钥    []byte

*/
type Wallet struct{
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte

}

func NewWallet()*Wallet {
	privateKey,publicKey:= newKeyPair()
	var wallet =&Wallet{
		PrivateKey: privateKey,
		PublicKey: publicKey,
	}
	return wallet
}

func newKeyPair()(ecdsa.PrivateKey,[]byte){
	/*
		1、获取一个椭圆
		2、通过椭圆曲线算法生成私钥
		3、通过私钥生成公钥
	*/
	curve:=elliptic.P256()
	privateKey,err:=ecdsa.GenerateKey(curve,rand.Reader)
	if err!=nil{
		log.Panicf("Generate Private Key Error: %v\n",err)
	}
	publicKey:=append(privateKey.PublicKey.X.Bytes(),privateKey.PublicKey.Y.Bytes()...)
	return *privateKey,publicKey
}


func (w *Wallet)GenerateAddress()[]byte{
	publicKeyHash:= GeneratePublicKeyHash(w.PublicKey)
	checkSum:= generateCheckSum(publicKeyHash)
	var data = [][]byte{
		publicKeyHash,
		checkSum,
	}
	input:=bytes.Join(data,[]byte{})
	address:=crypto.Base58Encode(input)
	return address
}

// GeneratePublicKeyHash : 生成公钥哈希
func GeneratePublicKeyHash(publicKey []byte)[]byte{
	/*
		生成PublicKeyHash的过程包括：
			1、对publicKey进行sha256运算
			2、对上述结果进行ripemd160运算
	*/
	sha:=sha256.New()
	sha.Write(publicKey)
	hash:=sha.Sum(nil)
	ripemd:=ripemd160.New()
	ripemd.Write(hash)
	publicKeyHash:=ripemd.Sum(nil)
	return publicKeyHash
}

// generateCheckSum : 生成校验和
func generateCheckSum(publicKeyHash []byte)[]byte{
	/*
		生成checkSum的过程：
			对publicKeyHash进行两次sha256运算
	*/
	hash:=sha256.Sum256(publicKeyHash)
	checkSum:=sha256.Sum256(hash[:])
	return checkSum[:]
}

func IsValidAddress(address []byte)bool{
	/*
		验证地址的有效性：
			1、对地址进行base58解码
			2、对解码后的字节数组进行拆分
			3、获取PublicKeyHash
			4、获取校验和
			5、将PublicKeyHash进行双哈希运算，比较校验和
	*/
	var checkSumLen=32
	decodeBytes:=crypto.Base58Decode(address)
	publicKeyHash:=decodeBytes[:len(decodeBytes)-checkSumLen]
	checkSumBytes:=decodeBytes[len(decodeBytes)-checkSumLen:]
	toVerifyCheckSum:= generateCheckSum(publicKeyHash)
	if bytes.Equal(checkSumBytes,toVerifyCheckSum){
		return true
	}
	return false
}