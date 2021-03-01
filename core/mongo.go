package core

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

const (
	MongoDBURI="mongodb://localhost:27017"
	DataBaseName="BlockChain"
	BlockTableName="Block"
	TipTableName="Tip"
	UTXOTableName="UTXO"

)

// ConnectMongo : 连接mongo数据库
func ConnectMongo() *mongo.Client{
	clientOptions := options.Client().ApplyURI(MongoDBURI)
	client,err:=mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		log.Fatal(err)
	}
	err=client.Ping(context.TODO(),nil)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("connect successfully")
	return client
}

// DisconnectMongo : 关闭mongo数据库连接
func DisconnectMongo(client *mongo.Client){
	err:=client.Disconnect(context.TODO())
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Disconnect successfully")
}

// InsertSensor : 插入文档
func InsertSensor(client *mongo.Client,document interface{})(insertID primitive.ObjectID){

	collection:=client.Database(DataBaseName).Collection(BlockTableName)
	insertResult,err:=collection.InsertOne(context.TODO(),document)
	if err!=nil{
		log.Fatal(err)
	}
	if insertResult!=nil{
		insertID=insertResult.InsertedID.(primitive.ObjectID)
	}
	fmt.Println("Insert block successfully")
	return insertID
}

// InsertSensors : 插入多条文档
func InsertSensors(client *mongo.Client,documents []interface{})(insertIDs []interface{}){
	collection:=client.Database(DataBaseName).Collection(BlockTableName)
	insertResult,err:=collection.InsertMany(context.TODO(),documents)
	if err!=nil{
		log.Fatal(err)
	}
	if insertResult!=nil{
		insertIDs=insertResult.InsertedIDs
	}
	return insertIDs
}

// QuerySensor : 查询文档
func QuerySensor(client *mongo.Client,key interface{})*Block {
	filter:=bson.M{"curblockhash":key}
	collection:=client.Database(DataBaseName).Collection(BlockTableName)
	result:=collection.FindOne(context.TODO(),filter)
	if result!=nil && result.Err()!=nil{
		log.Panicf("Query mongo Error:%v",result.Err().Error())
	}
	var block =&Block{}
	err:=result.Decode(block)
	if err!=nil{
		fmt.Println(err)
	}
	return block
}

// UpdateTip : 更新Tip
func UpdateTip(client *mongo.Client,tip interface{}){
	collection:=client.Database(DataBaseName).Collection(TipTableName)
	filter:=bson.D{}
	value:=bson.M{"l":tip}
	_ = collection.FindOneAndUpdate(context.TODO(),filter,bson.M{"$set":value})
	//if singleResult!=nil{
	//	fmt.Println(singleResult.Decode(&Block{}))
	//}
	fmt.Println("Update tip successfully")
}

// InsertTip : 插入tip
func InsertTip(client *mongo.Client,tip interface{})(insertID primitive.ObjectID){
	collection:=client.Database(DataBaseName).Collection(TipTableName)
	insertResult,err:=collection.InsertOne(context.TODO(),tip)
	if err!=nil{
		fmt.Println(err)
	}
	if insertResult!=nil{
		insertID=insertResult.InsertedID.(primitive.ObjectID)
	}
	fmt.Println("Insert tip successfully")
	return insertID
}

// QueryTip : 查询Tip
func QueryTip(client *mongo.Client)*Tip {
	filter:=bson.D{}
	collection:=client.Database(DataBaseName).Collection(TipTableName)
	cursor , err :=collection.Find(context.TODO(),filter)
	if err!=nil{
		fmt.Println(err)
	}
	defer cursor.Close(context.TODO())
	var tip =&Tip{}
	for cursor.Next(context.TODO()){
		err = cursor.Decode(tip)
		if err!=nil{
			fmt.Println(err.Error())
		}
	}
	return tip
}

// ResetUTXOTable : 重置UTXO Table
func ResetUTXOTable(blockchain *BlockChain){
	client:=blockchain.Client
	// 删除UTXO集合 : 因为在区块链创建时调用，所以先删除所有的UTXO集合
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	err:=collection.Drop(context.TODO())
	if err!=nil{
		log.Panicf("Drop UTXO Table Error:%v\n",err)
	}
	// 创建UTXO集合 : 存储UTXO集合，来提高原来传统查询的效率
	utxos:=blockchain.FindAllUTXO()
	for keyHash,outputs:= range utxos{
		txOutputs:= TxOutputs{keyHash,outputs}
		_, err := collection.InsertOne(context.TODO(), txOutputs)
		if err!=nil{
			log.Panicf("Insert UTXO ERROR:%v\n",err)
		}
	}
}

// FindUTXOTable : 查询UTXO集
func FindUTXOTable(blockchain *BlockChain,address string)[]*UTXO {
	var utxos []*UTXO
	client:=blockchain.Client
	filter:=bson.D{}
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	cursor , err :=collection.Find(context.TODO(),filter)
	if err!=nil{
		log.Printf("Query UTXO ERROR:%V\n",err)
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()){
		var txOutputs =&TxOutputs{}
		err = cursor.Decode(&txOutputs)
		if err!=nil{
			fmt.Println(err.Error())
		}
		for _,utxo:=range txOutputs.TxOutputs{
			if utxo.CheckPublicKeyWithAddress(address){
				utxoSingle:= UTXO{
					Output: utxo,
				}
				utxos=append(utxos, &utxoSingle)
			}
		}
	}
	return utxos
}

//
func QueryTxOutputs(client *mongo.Client,txHash string)*TxOutputs {
	filter:=bson.M{"txhash":txHash}
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	result:=collection.FindOne(context.TODO(),filter)
	if result==nil && result.Err()!=nil{
		log.Panicf("Query TxOutputs ERROR:%v",result.Err().Error())
	}
	var txOutputs = TxOutputs{}
	err:=result.Decode(&txOutputs)
	if err!=nil{
		fmt.Println(err)
	}
	return &txOutputs
}
//
func DeleteUTXOTable(client *mongo.Client,txHash string)*mongo.DeleteResult{
	filter:=bson.M{"txhash":txHash}
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	deleteResult,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		log.Panicf("Delete TxOutputs ERROR:%v",err)
	}
	return deleteResult
}

//
func UpdateUTXOTable(client *mongo.Client,txHash string,txOutputs []*TxOutput){
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	filter:=bson.M{"txhash":txHash}
	value:=bson.M{"txoutputs":txOutputs}
	_ = collection.FindOneAndUpdate(context.TODO(),filter,bson.M{"$set":value})
	//if singleResult!=nil{
	//	fmt.Println(singleResult.Decode(&Block{}))
	//}
	fmt.Println("Update UTXO Table successfully")
}

func InsertUTXOTable(client *mongo.Client,txHash string,txOutputs []*TxOutput){
	collection:=client.Database(DataBaseName).Collection(UTXOTableName)
	_, err := collection.InsertOne(context.TODO(), TxOutputs{txHash,txOutputs})
	if err!=nil{
		log.Panicf("Insert UTXO ERROR:%v\n",err)
	}
}