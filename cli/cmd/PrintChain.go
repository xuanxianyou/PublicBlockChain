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
	"os"
)

// PrintChainCmd represents the PrintChain command
var PrintChainCmd = &cobra.Command{
	Use:   "PrintChain",
	Short: "Travel the whole BlockChain",
	Long: `Travel the whole BlockChain`,
	Run: func(cmd *cobra.Command, args []string) {
		client:= core.ConnectMongo()
		defer core.DisconnectMongo(client)
		blockchain:=core.GetBlockChain(client)
		if blockchain.IsGenesisBlockExisted(){
			blockchain.PrintBlockChain()
		}else{
			fmt.Println("The Genesis Block isn`t Existed...")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(PrintChainCmd)

}
