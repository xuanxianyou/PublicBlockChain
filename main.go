package main

import (
	"GoProject/BlockChain/PublicBlockChain/core"
)

func main() {
	client:= core.ConnectMongo()
	defer core.DisconnectMongo(client)
	core.QueryTxOutputs(client,"f706c8f001a66998a87146d4757fd5104496fe43972c9c9fd65ac7a52d9f8923")
}