/*
Copyright © 2021 Kaneziki <1848224883@qq.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"GoProject/BlockChain/PublicBlockChain/core"
	"fmt"

	"github.com/spf13/cobra"
)

// QueryBalanceCmd represents the QueryBalance command
var QueryBalanceCmd = &cobra.Command{
	Use:   "QueryBalance",
	Short: "Query the Balance according to the address",
	Long: `Query the Balance according to the address`,
	Run: func(cmd *cobra.Command, args []string) {
		var balance int
		address,err:=cmd.Flags().GetString("address")
		if err!=nil{
			fmt.Println(err)
		}
		client:= core.ConnectMongo()
		defer core.DisconnectMongo(client)
		blockchain:=core.GetBlockChain(client)
		//utxos:=blockchain.FindUnspentUTXOs(address,[]*core.Transaction{})
		//for _,utxo:=range utxos{
		//	balance+=utxo.Output.Value
		//}
		utxoSet:= core.UTXOSet{
			Blockchain: blockchain,
		}
		balance=utxoSet.GetBalance(address)
		fmt.Printf("%s`s balance is %d Hacoin\n",address,balance)
	},
}

func init() {
	rootCmd.AddCommand(QueryBalanceCmd)
	QueryBalanceCmd.Flags().StringP("address","a","","Specify user address")
}
