package core

import "crypto/sha256"

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

//
func NewMerkleTree(txHashes [][]byte)*MerkleTree{
	// 结点列表
	var nodes =[]*MerkleNode{}
	// 判断结点个数
	if len(txHashes)%2==1{
		// 如果有奇数个结点
		txHashes = append(txHashes, txHashes[len(txHashes)-1])
	}
	for _,data:=range txHashes{
		// 遍历所有交易哈希，生成叶子结点
		node:=NewMerkleNode(nil,nil,data)
		nodes=append(nodes, node)
	}
	for i:=0;i<len(txHashes);i++{
		var parentNodes =[]*MerkleNode{}
		for j:=0;j<len(nodes);j+=2{
			parentNode:=NewMerkleNode(nodes[j],nodes[j+1],nil)
			parentNodes=append(parentNodes, parentNode)
		}
		if len(parentNodes)%2==1{
			parentNodes=append(parentNodes, parentNodes[len(parentNodes)-1])
		}
		nodes=parentNodes
	}
	merkleTree:=&MerkleTree{nodes[0]}
	return merkleTree
}

//
func NewMerkleNode(left *MerkleNode,right *MerkleNode,data []byte)*MerkleNode{
	var node = &MerkleNode{}
	if left==nil && right==nil{
		// 叶子结点
		hash:=sha256.Sum256(data)
		node.Data=hash[:]
	}else{
		// 非叶子结点
		preHash:=append(left.Data,right.Data...)
		hash:=sha256.Sum256(preHash)
		node.Data=hash[:]

	}
	node.Left=left
	node.Right=right
	return node
}