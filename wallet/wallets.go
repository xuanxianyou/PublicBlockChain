package wallet

import (
	"bytes"
	"crypto/elliptic"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

const(
	FilePath="Wallets.dat"
)

/*
	Wallets : 钱包的集合
		一个地址对应一个钱包
		address---->wallet
*/
type Wallets struct{
	WalletsMap map[string]*Wallet
}

// NewWallets : 生成钱包集合
func NewWallets()*Wallets {
	// 从文件中对钱包集合初始化
	var wallets Wallets
	if _,err:=ioutil.ReadFile(FilePath);os.IsNotExist(err){
		// 如果文件不存在，创建新的wallet
		wallets.WalletsMap=make(map[string]*Wallet)
		return &wallets
	}
	//如果文件存在，从文件中载入
	walletsBytes:= loadWallets()
	if len(walletsBytes)==0{
		// 如果文件存在，内容为空
		wallets.WalletsMap=make(map[string]*Wallet)
		return &wallets
	}
	wallets = *deSerialize(walletsBytes)
	return &wallets
}

// CreateWallet : 创建钱包并加入钱包集合
func (w *Wallets)CreateWallet(){
	wallet:= NewWallet()
	address:=wallet.GenerateAddress()
	fmt.Printf("Your address:%s\n",address)
	w.WalletsMap[string(address)]=wallet
	// 本地持久化
	w.saveWallets()
}

// ShowAddress : 查看钱包地址
func (w *Wallets)ShowAddresses(){
	for address,_:=range w.WalletsMap{
		fmt.Println(address)
	}
}

// saveWallets : 将钱包存储到文件中
func (w *Wallets)saveWallets(){
	walletsBytes:=w.serialize()
	err:=ioutil.WriteFile(FilePath,walletsBytes,0644)
	if err!=nil{
		log.Panicf("Write wallets to file Error:%v\n",err)
	}
}

// loadWallets : 从文件中读取wallets
func loadWallets()[]byte{
	walletsBytes,err:=ioutil.ReadFile(FilePath)
	if err!=nil{
		log.Panicf("Read wallets from file Error:%v\n",err)
	}
	return walletsBytes
}

// serialize : 钱包集合序列化
func (w *Wallets)serialize()[]byte{
	var buffer bytes.Buffer
	gob.Register(elliptic.P256())
	encoder:=gob.NewEncoder(&buffer)
	err:=encoder.Encode(&w)
	if err!=nil{
		log.Panicf("wallets serialize Error:%v\n",err)
	}
	return buffer.Bytes()
}

// deSerialize : 钱包集合反序列化
func deSerialize(walletsBytes []byte)*Wallets {
	var wallets Wallets
	gob.Register(elliptic.P256())
	decoder:=gob.NewDecoder(bytes.NewReader(walletsBytes))
	if err:=decoder.Decode(&wallets);err!=nil{
		log.Panicf("wallets deserialize Error:%v\n",err)
	}
	return &wallets
}