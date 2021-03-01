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
	"os"

	"github.com/spf13/cobra"
)

// TransferAccountCmd represents the TransferAccount command
var TransferAccountCmd = &cobra.Command{
	Use:   "TransferAccount",
	Short: "Transfer money From your Account To another Account",
	Long: `Transfer money From your Account To another Account`,
	Run: func(cmd *cobra.Command, args []string) {
		fromArray,err:=cmd.Flags().GetStringArray("from")
		toArray,err:=cmd.Flags().GetStringArray("to")
		moneyArray,err:=cmd.Flags().GetIntSlice("money")
		if err!=nil{
			fmt.Println(err)
			os.Exit(1)
		}
		// 错误检测
		if len(fromArray)==0 || len(toArray)==0 ||len(moneyArray)==0{
			// 参数内容为空检测
			fmt.Println("The number of args can`t be zero!")
			os.Exit(1)
		}
		if len(fromArray)!=len(toArray) && len(fromArray)!=len(moneyArray){
			// 参数数目不等检测
			fmt.Println("The number of args isn`t match!")
			os.Exit(1)
		}
		for i:=0;i<len(fromArray);i++{
			// 非法参数检测
			if fromArray[i]==toArray[i]{
				fmt.Printf("The address in the %d transaction is invalid",i)
			}
			if moneyArray[i]==0{
				fmt.Printf("The money in the %d transaction is invalid",i)
			}
		}
		client:= core.ConnectMongo()
		defer core.DisconnectMongo(client)
		blockchain:=core.GetBlockChain(client)
		if blockchain.IsGenesisBlockExisted(){
			var transactions []*core.Transaction
			for i:=0;i<len(fromArray);i++{
				from :=fromArray[i]
				to := toArray[i]
				money:=moneyArray[i]
				transaction:= core.NewTransferTransaction(from,to,money,blockchain,transactions)
				transactions=append(transactions, transaction)
			}
			blockchain.AddBlock(transactions)
			// 更新UTXO table
			utxoSet:=&core.UTXOSet{Blockchain: blockchain}
			utxoSet.UpdateUTXOTable()
		}else{
			fmt.Println("The Genesis Block isn`t Existed")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(TransferAccountCmd)
	TransferAccountCmd.Flags().StringArrayP("from","f",[]string{},"Original address")
	TransferAccountCmd.Flags().StringArrayP("to","t",[]string{},"Target address")
	TransferAccountCmd.Flags().IntSliceP("money","m",[]int{},"money account")

}
