/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"GoProject/BlockChain/PublicBlockChain/wallet"
	"fmt"
	"github.com/spf13/cobra"
)

// CreateAddressCmd represents the CreateAddress command
var CreateAddressCmd = &cobra.Command{
	Use:   "CreateAddress",
	Short: "Create your account address",
	Long: `Create your account address`,
	Run: func(cmd *cobra.Command, args []string) {
		var wallets = wallet.NewWallets()
		if wallets!=nil{
			wallets.CreateWallet()
		}else{
			fmt.Println("Your address is null!")
		}
	},
}

func init() {
	rootCmd.AddCommand(CreateAddressCmd)
}
