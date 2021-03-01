package core

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Iterator struct {
	Client       *mongo.Client
	CurrentBlock *Block
}

// NewIterator : 创建迭代器
func NewIterator(client *mongo.Client,tip []byte) *Iterator {
	var iterator *Iterator
	block:= QuerySensor(client,tip)
	iterator = &Iterator{
		Client: client,
		CurrentBlock: block,
	}
	return iterator
}

// HasNext : 判断是否还有下一个区块
func (i *Iterator) HasNext() bool {
	if i.CurrentBlock.Height==-1{
		return false
	}
	return true
}

// Next : 获取下一个区块
func (i *Iterator) Next() *Block {
	i.CurrentBlock= QuerySensor(i.Client,i.CurrentBlock.PreBlockHash)
	return i.CurrentBlock
}

// GetCurrentBlock :  获取最新的区块
func (i *Iterator)GetCurrentBlock()*Block{
	return i.CurrentBlock
}
