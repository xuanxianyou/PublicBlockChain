package main

import (
	"GoProject/BlockChain/PublicBlockChain/core"
	"flag"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"strconv"
)

type Cli struct{
	Client *mongo.Client
}

// PrintUsage : 打印HaHa的用法
func PrintUsage(){
	fmt.Println("Usage: HaHa [Flags] [Command]\n")
	//简要
	fmt.Println("A Client For HaHa Public BlockChain\n")
	//命令
	fmt.Println("Commands:")
	fmt.Println("\tCreateChain\tCreate a BlockChain\n")
	fmt.Println("\tAddBlock\tAdd a Block to BlockChain\n")
	fmt.Println("\tTransferAccount\tTransfer money From your Account To another Account\n")
	fmt.Println("\tPrintChain\tTravel the whole BlockChain\n")
	fmt.Println("\tQueryBalance\tQuery the Balance according to the address")

}

// NewCli : 生成Cli对象
func NewCli(client *mongo.Client)*Cli{
	return &Cli{
		Client: client,
	}
}

//IsVailed : 判断命令是否有效
func IsVailed()bool{
	if os.Args[0]!="HaHa"{
		return false
	}
	return true
}


func (c *Cli)CreateChain(address string){
	client:= core.ConnectMongo()
	defer core.DisconnectMongo(client)
	core.NewBlockChain(client,address)
}

func (c *Cli)AddBlock(transaction []*core.Transaction){
	blockchain:=core.GetBlockChain(c.Client)
	if blockchain.IsGenesisBlockExisted(){
		blockchain.AddBlock(transaction)
	}else{
		fmt.Println("The Genesis Block isn`t Existed")
		os.Exit(1)
	}
}

func (c *Cli)TransferAccount(from string,to string,money int){
	blockchain:=core.GetBlockChain(c.Client)
	if blockchain.IsGenesisBlockExisted(){
		var transactions []*core.Transaction
		transaction:= core.NewTransferTransaction(from,to,money,blockchain,nil)
		transactions=append(transactions, transaction)
		blockchain.AddBlock(transactions)
	}else{
		fmt.Println("The Genesis Block isn`t Existed")
		os.Exit(1)
	}
}

func (c *Cli)PrintChain(){
	blockchain:=core.GetBlockChain(c.Client)
	if blockchain.IsGenesisBlockExisted(){
		blockchain.PrintBlockChain()
	}else{
		fmt.Println("The Genesis Block isn`t Existed")
		os.Exit(1)
	}
}

func (c *Cli)QueryBalance(address string){
	var balance int
	blockchain:=core.GetBlockChain(c.Client)
	utxos:=blockchain.FindUnspentUTXOs(address)
	for _,utxo:=range utxos{
		balance+=utxo.Output.Value
	}
	fmt.Printf("%s`s balance is %d Hacoin\n",address,balance)
}

// Run : 运行Cli命令
func (c *Cli)Run(){
	if !IsVailed(){
		PrintUsage()
		os.Exit(1)
	}
	//创建区块链命令
	createChainCmd:=flag.NewFlagSet("CreateChain",flag.ExitOnError)
	address:=createChainCmd.String("address","","The HaHa BlockChain Creator`s address")
	//添加区块命令
	addBlockCmd:=flag.NewFlagSet("AddBlock",flag.ExitOnError)
	transaction:=addBlockCmd.String("data","","Add the Specified Data to the HaHa Blockchain")
	//转账命令
	transferAccountCmd:=flag.NewFlagSet("TransferAccount",flag.ExitOnError)
	from:=transferAccountCmd.String("from","","original Account")
	to:=transferAccountCmd.String("to","","target Account")
	money:=transferAccountCmd.String("money","","transfer money amount")
	//打印区块命令
	printChainCmd:=flag.NewFlagSet("PrintChain",flag.ExitOnError)
	//查询余额命令
	queryBalanceCmd:=flag.NewFlagSet("QueryBalance",flag.ExitOnError)
	queryAddress:=queryBalanceCmd.String("address","","Query the Balance according to the address")
	switch os.Args[len(os.Args)-1] {
	case "CreateChain":
		if err:=createChainCmd.Parse(os.Args[1:len(os.Args)-1]);err!=nil{
			log.Panic(err)
		}
	case "AddBlock":
		if err:=addBlockCmd.Parse(os.Args[1:len(os.Args)-1]);err!=nil{
			log.Panic(err)
		}
	case "TransferAccount":
		if err:=transferAccountCmd.Parse(os.Args[1:len(os.Args)-1]);err!=nil{
			log.Panic(err)
		}
	case "PrintChain":
		if err:=printChainCmd.Parse(os.Args[1:len(os.Args)-1]);err!=nil{
			log.Panic(err)
		}
	case "QueryBalance":
		if err:=queryBalanceCmd.Parse(os.Args[1:len(os.Args)-1]);err!=nil{
			log.Panic(err)
		}
	default:
		PrintUsage()
		os.Exit(1)
	}
	if createChainCmd.Parsed(){
		if *address==""{
			PrintUsage()
			os.Exit(1)
		}else{
			c.CreateChain(*address)
		}
	}
	if addBlockCmd.Parsed(){
		if *transaction==""{
			PrintUsage()
			os.Exit(1)
		}else{
			c.AddBlock([]*transaction.Transaction{})
		}
	}
	if transferAccountCmd.Parsed(){
		if *from==""{
			fmt.Println("Invalid original address!")
			PrintUsage()
			os.Exit(1)
		}
		if *to==""{
			fmt.Println("Invalid target address!")
			PrintUsage()
			os.Exit(1)
		}
		if *money==""{
			fmt.Println("Invalid args money")
			PrintUsage()
			os.Exit(1)
		}
		moneyInt,err:=strconv.Atoi(*money)
		if err!=nil{
			log.Panic(err)
		}
		c.TransferAccount(*from,*to,moneyInt)


	}
	if printChainCmd.Parsed(){
		c.PrintChain()
	}
	if queryBalanceCmd.Parsed(){
		if *queryAddress==""{
			PrintUsage()
			os.Exit(1)
		}else{
			c.QueryBalance(*queryAddress)
		}
	}

}
