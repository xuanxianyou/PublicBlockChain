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
	"log"
)


// CreateChainCmd represents the CreateChain command
var CreateChainCmd = &cobra.Command{
	Use:   "CreateChain",
	Short: "Create HaHa BlockChain",
	Long: `Create HaHa BlockChain`,
	Run: func(cmd *cobra.Command, args []string) {
		// 连接数据库
		client:= core.ConnectMongo()
		defer core.DisconnectMongo(client)
		blockchain:=core.GetBlockChain(client)
		if blockchain.IsGenesisBlockExisted(){
			//如果创世区块已经存在
			fmt.Println("The HaHa Blockchain is Existed...")
		}else{
			address,err:=cmd.Flags().GetString("address")
			if err!=nil{
				log.Println(err)
			}
			if address==""{
				fmt.Println("Invalid address...")
			}else{
				blockchain=core.NewBlockChain(client,address)

				utxoSet:= core.UTXOSet{
					Blockchain: blockchain,
				}
				utxoSet.ResetUTXOSet()
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(CreateChainCmd)
	CreateChainCmd.Flags().StringP("address","a","","Specify user address")
}
