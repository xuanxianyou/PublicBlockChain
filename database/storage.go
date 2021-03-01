package database

import (
	"GoProject/BlockChain/PublicBlockChain/core"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
)

//StorageBlock : 添加Block数据
func StorageBlock(key []byte,value []byte){
	db,err:=bolt.Open(core.DBPath,0600,nil)
	if err!=nil{
		log.Panic(err.Error())
	}
	if db!=nil{
		defer db.Close()
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(core.TableName))
			if b == nil {
				b, err = tx.CreateBucket([]byte(core.TableName))
				if err != nil {
					log.Panic(err.Error())
				}
			}
			err=b.Put(key,value)
			if err!=nil{
				log.Panic(err.Error())
			}
			return nil
		})
		if err!=nil{
			log.Panic(err.Error())
		}
	}
}

//ReadBlock : 返回上一块区块
func ReadBlock(key []byte)*core.Block {
	var block *core.Block
	db,err:=bolt.Open(core.DBPath,0600,nil)
	if err!=nil{
		log.Panic(err.Error())
	}
	if db!=nil{
		defer db.Close()
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(core.TableName))
			if b == nil {
				b, err := tx.CreateBucket([]byte(core.TableName))
				if err != nil {
					log.Panic(err.Error())
				}
				blockBytes:=b.Get(key)
				block= core.Deserialize(blockBytes)
			}
			return nil
		})
		if err!=nil{
			log.Panic(err.Error())
		}
	}
	return block
}

//SetTip : 设置Tip
func SetTip(tip []byte){
	db,err:=bolt.Open(core.DBPath,0600,nil)
	if err!=nil{
		log.Panic(err.Error())
	}
	if db!=nil{
		defer db.Close()
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(core.TableName))
			if b == nil {
				b, err = tx.CreateBucket([]byte(core.TableName))
				if err != nil {
					log.Panic(err.Error())
				}
			}
			err=b.Put([]byte("tip"),tip)
			if err!=nil{
				log.Panic(err.Error())
			}
			return nil
		})
		if err!=nil{
			log.Panic(err.Error())
		}
	}
}

//GetTip : 获取tip的值
func GetTip(tip []byte) {
	db,err:=bolt.Open(core.DBPath,0600,nil)
	if err!=nil{
		log.Panic(err.Error())
	}
	if db!=nil{
		defer db.Close()
		_ = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(core.TableName))
			for _,v := range b.Get([]byte("tip")){
				tip=append(tip, v)
			}
			return nil
		})
	}
	fmt.Println(tip)
}