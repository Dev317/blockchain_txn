package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/Dev317/go_blockchain/blockchain/eth"
	"github.com/Dev317/go_blockchain/config"
	"github.com/spf13/viper"
)

func main () {

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		slog.Error("Error reading config file", err)
		os.Exit(1)
	}

	var config config.Config
	viper.Unmarshal(&config)

	client, err := ethclient.Dial(config.Ethereum.RpcURL)
	if err != nil {
		slog.Error("Error connecting to Ethereum testnet", err)
		os.Exit(1)
	}

	slog.Info("Establish a connection to the Ethereum testnet")

	// from_address := utils.GenerateKey()
	// to_address := utils.GenerateKey()

	// from_account := common.HexToAddress(from_address)
	// to_account := common.HexToAddress(to_address)

	from_account := common.HexToAddress(config.Ethereum.FromAccount.Address)
	to_account := common.HexToAddress(config.Ethereum.ToAccount.Address)

    from_balance, err := client.BalanceAt(context.Background(), from_account, nil)
    if err != nil {
        slog.Error("Error getting balance", err)
    }

    slog.Info(fmt.Sprintf("Address %s has balance: %s wei", from_account.String(), from_balance.String()))

	to_balance, err := client.BalanceAt(context.Background(), to_account, nil)
    if err != nil {
        slog.Error("Error getting balance", err)
    }

    slog.Info(fmt.Sprintf("Address %s has balance: %s wei", to_account.String(), to_balance.String()))

	amount := big.NewInt(1000000)

	slog.Debug("Sending transaction...")
	utils.Send(config.Ethereum.FromAccount.PrivateKey,
			   config.Ethereum.ToAccount.Address,
			   amount,
			   client)
	slog.Info("Transaction sent successfully")

	time.Sleep(30 * time.Second)
	from_balance, err = client.BalanceAt(context.Background(), from_account, nil)
    if err != nil {
        slog.Error("Error getting balance", err)
    }

    slog.Info(fmt.Sprintf("Address %s has balance: %s wei", from_account.String(), from_balance.String()))

	to_balance, err = client.BalanceAt(context.Background(), to_account, nil)
    if err != nil {
        slog.Error("Error getting balance", err)
    }

    slog.Info(fmt.Sprintf("Address %s has balance: %s wei", to_account.String(), to_balance.String()))
}