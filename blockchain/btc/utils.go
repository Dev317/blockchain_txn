package btcutils

import (
	"bytes"
	"encoding/hex"
	// "encoding/json"
	"fmt"
	// "io"
	// "net/http"
	"time"
	"log/slog"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)


type BlockChairResp struct {
	Address             string         `json:"address"`
	TotalReceived       int64          `json:"total_received"`
	TotalSent           int64          `json:"total_sent"`
	Balance             int64          `json:"balance"`
	UnconfirmedBalance  int64          `json:"unconfirmed_balance"`
	FinalBalance        int64          `json:"final_balance"`
	NTx                 int            `json:"n_tx"`
	UnconfirmedNTx      int            `json:"unconfirmed_n_tx"`
	FinalNTx            int            `json:"final_n_tx"`
	Txs                 []Transaction  `json:"txs"`
}

// Transaction represents each transaction in the JSON response.
type Transaction struct {
	BlockHash      string         `json:"block_hash"`
	BlockHeight    int64          `json:"block_height"`
	BlockIndex     int            `json:"block_index"`
	Hash           string         `json:"hash"`
	Addresses      []string       `json:"addresses"`
	Total          int64          `json:"total"`
	Fees           int64          `json:"fees"`
	Size           int            `json:"size"`
	VSize          int            `json:"vsize"`
	Preference     string         `json:"preference"`
	RelayedBy      string         `json:"relayed_by"`
	Confirmed      time.Time      `json:"confirmed"`
	Received       time.Time      `json:"received"`
	Ver            int            `json:"ver"`
	LockTime       int64          `json:"lock_time"`
	DoubleSpend    bool           `json:"double_spend"`
	VinSz          int            `json:"vin_sz"`
	VoutSz         int            `json:"vout_sz"`
	OptInRBF       bool           `json:"opt_in_rbf"`
	Confirmations  int            `json:"confirmations"`
	Confidence     float64        `json:"confidence"`
	Inputs         []Input        `json:"inputs"`
	Outputs        []Output       `json:"outputs"`
}

// Input represents each input in a transaction.
type Input struct {
	PrevHash     string   `json:"prev_hash"`
	OutputIndex  int      `json:"output_index"`
	OutputValue  int64    `json:"output_value"`
	Sequence     int64    `json:"sequence"`
	Addresses    []string `json:"addresses"`
	ScriptType   string   `json:"script_type"`
	Age          int64    `json:"age"`
	Witness      []string `json:"witness"`
}

// Output represents each output in a transaction.
type Output struct {
	Value       int64    `json:"value"`
	Script      string   `json:"script"`
	Addresses   []string `json:"addresses"`
	ScriptType  string   `json:"script_type"`
	SpentBy     string   `json:"spent_by,omitempty"` // Use omitempty for optional fields
}

func NewTx() (*wire.MsgTx, error) {
	return wire.NewMsgTx(wire.TxVersion), nil
}

func GetUTXO(address string) (string, int64, string, error) {

	// Provide your url to get UTXOs, read the response
	// unmarshal it, and extract necessary data
	// newURL := fmt.Sprintf("https://api.blockcypher.com/v1/btc/test3/addrs/%s/full?limit=5", "tb1qcxkkj7pter9dv3mdsgjv02dzyxj73p6qfp3cdws5g75qdgte7zashdg2ql")

	// response, err := http.Get(newURL)
	// if err != nil {
	// fmt.Println("error in GetUTXO, http.Get")
	// return "", 0, "", err
	// }
	// defer response.Body.Close()
	// body, err := io.ReadAll(response.Body)

	// var blockChairResp = BlockChairResp{}
	// err = json.Unmarshal(body, &blockChairResp)
	// if err != nil {
	// 	fmt.Println("error in GetUTXO, json.Unmarshal")
	// 	return  "", 0, "", err
	// }

	// fmt.Print("UTXO latest txn: ", blockChairResp.Txs[0])


	var previousTxid string = "9d37c847cbc9ab76ec6ccdbfe84833f08df179d2da04fa3ae212475aa94ef809"
	var balance int64 = 30000
	var pubKeyScript string = "0020c1ad69782bc8cad6476d8224c7a9a221a5e88740486386ba1447a806a179f0bb"
	return previousTxid, balance, pubKeyScript, nil
}

func CreateTx(privKey string, destination string, amount int64) (string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	// use TestNet3Params for interacting with bitcoin testnet
	// if we want to interact with main net should use MainNetParams
	addrPubKey, err := btcutil.NewAddressPubKey(wif.PrivKey.PubKey().SerializeUncompressed(), &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	txid, balance, pkScript, err := GetUTXO(addrPubKey.EncodeAddress())
	if err != nil {
		return "", err
	}

	/*
	 * 1 or unit-amount in Bitcoin is equal to 1 satoshi and 1 Bitcoin = 100000000 satoshi
	 */

	// checking for sufficiency of account
	if balance < amount {
		return "", fmt.Errorf("the balance of the account is not sufficient")
	}

	// extracting destination address as []byte from function argument (destination string)
	destinationAddr, err := btcutil.DecodeAddress(destination, &chaincfg.TestNet3Params)
	if err != nil {
		return "", err
	}

	destinationAddrByte, err := txscript.PayToAddrScript(destinationAddr)
	if err != nil {
		return "", err
	}

	// creating a new bitcoin transaction, different sections of the tx, including
	// input list (contain UTXOs) and outputlist (contain destination address and usually our address)
	// in next steps, sections will be field and pass to sign
	redeemTx, err := NewTx()
	if err != nil {
		return "", err
	}

	utxoHash, err := chainhash.NewHashFromStr(txid)
	if err != nil {
		return "", err
	}

	// the second argument is vout or Tx-index, which is the index
	// of spending UTXO in the transaction that Txid referred to
	// in this case is 1, but can vary different numbers
	outPoint := wire.NewOutPoint(utxoHash, 1)

	// making the input, and adding it to transaction
	txIn := wire.NewTxIn(outPoint, nil, nil)
	redeemTx.AddTxIn(txIn)

	// adding the destination address and the amount to
	// the transaction as output
	redeemTxOut := wire.NewTxOut(amount, destinationAddrByte)
	redeemTx.AddTxOut(redeemTxOut)

	// now sign the transaction
	finalRawTx, _ := SignTx(privKey, pkScript, redeemTx)

	return finalRawTx, nil
}

func SignTx(privKey string, pkScript string, redeemTx *wire.MsgTx) (string, error) {

	wif, err := btcutil.DecodeWIF(privKey)
	if err != nil {
		return "", err
	}

	sourcePKScript, err := hex.DecodeString(pkScript)
	if err != nil {
		return "", nil
	}

	// since there is only one input in our transaction
	// we use 0 as second argument, if the transaction
	// has more args, should pass related index
	signature, err := txscript.SignatureScript(redeemTx, 0, sourcePKScript, txscript.SigHashAll, wif.PrivKey, false)
	if err != nil {
		return "", nil
	}

	// since there is only one input, and want to add
	// signature to it use 0 as index
	redeemTx.TxIn[0].SignatureScript = signature

	var signedTx bytes.Buffer
	redeemTx.Serialize(&signedTx)

	hexSignedTx := hex.EncodeToString(signedTx.Bytes())

	return hexSignedTx, nil
}


func SendTxn(prevTxHash string, outputIdx uint, wif string, destAddr string, amount int64) {
	outputIndex := uint32(outputIdx)

	// Decode WIF
	wifObj, err := btcutil.DecodeWIF(wif)
	if err != nil {
		slog.Error("Error decoding WIF: %v", err)
	}

	// Parse the destination address
	addr, err := btcutil.DecodeAddress(destAddr, &chaincfg.TestNet3Params)
	if err != nil {
		slog.Error("Error decoding address: %v", err)
	}

	// Create a new transaction
	tx := wire.NewMsgTx(wire.TxVersion)

	// Add the input (reference the previous transaction)
	hash, err := chainhash.NewHashFromStr(prevTxHash)
	if err != nil {
		slog.Error("Error parsing hash: %v", err)
	}

	outPoint := wire.NewOutPoint(hash, outputIndex)
	txIn := wire.NewTxIn(outPoint, nil, nil)
	tx.AddTxIn(txIn)

	// Define the amount to send (in Satoshi). For simplicity, we're not handling change here.
	amountToSend := int64(amount)
	pkScript, err := txscript.PayToAddrScript(addr)
	if err != nil {
		slog.Error("Error generating pkScript: %v", err)
	}

	txOut := wire.NewTxOut(amountToSend, pkScript)
	tx.AddTxOut(txOut)

	// Sign the transaction
	signature, err := txscript.SignatureScript(tx, 0, pkScript, txscript.SigHashAll, wifObj.PrivKey, true)
	if err != nil {
		slog.Error("Error creating signature: %v", err)
	}

	// Set the signature as the witness for the transaction.
	// This is a critical step for SegWit transactions.
	tx.TxIn[0].Witness = wire.TxWitness{signature, wifObj.PrivKey.PubKey().SerializeCompressed()}


	// Serialize and display the transaction
	var signedTx bytes.Buffer
	if err := tx.Serialize(&signedTx); err != nil {
		slog.Error("Error serializing transaction: %v", err)
	}
	fmt.Printf("Signed Transaction: %x\n", signedTx.Bytes())

}
