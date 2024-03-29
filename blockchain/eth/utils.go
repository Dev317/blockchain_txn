package utils

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)


func GenerateKey() (string){
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		slog.Error("Error generating key", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	slog.Info(fmt.Sprintf("Generated private key: %s", hexutil.Encode(privateKeyBytes)[2:]))

	publicKey := privateKey.Public()

    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        slog.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	slog.Info(fmt.Sprintf("Generated public key: %s", hexutil.Encode(publicKeyBytes)[4:]))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	slog.Info(fmt.Sprintf("Generated address: %s", address))

	return address

}

func Send(pk string, toHexAddress string, value *big.Int, client *ethclient.Client) {
	privateKey, err := crypto.HexToECDSA(pk)
    if err != nil {
        slog.Error("Error converting hex to ECDSA", err)
    }

    publicKey := privateKey.Public()
    publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
    if !ok {
        slog.Error("Cannot assert type: publicKey is not of type *ecdsa.PublicKey")
    }

    fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
    nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
    if err != nil {
        slog.Error("",err)
    }

    gasLimit := uint64(21000)
    gasPrice, err := client.SuggestGasPrice(context.Background())
    if err != nil {
        slog.Error("",err)
    }

    toAddress := common.HexToAddress(toHexAddress)
    var data []byte
    tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

    chainID, err := client.NetworkID(context.Background())
    if err != nil {
        slog.Error("Error in getting network ID",err)
    }

    signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
    if err != nil {
        slog.Error("Error in signing transaction",err)
    }

    err = client.SendTransaction(context.Background(), signedTx)
    if err != nil {
        slog.Error("Error in sending transaction",err)
    }

    slog.Info(fmt.Sprintf("Tx sent: %s", signedTx.Hash().Hex()))
}