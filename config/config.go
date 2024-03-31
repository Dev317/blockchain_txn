package config

type Config struct {
    Ethereum Ethereum `mapstructure:"ethereum"`
    Bitcoin Bitcoin `mapstructure:"bitcoin"`
}

type Bitcoin struct {
    PrevTxHash string `mapstructure:"prev_tx_hash"`
    WIF string `mapstructure:"wif"`
    DestAddr string `mapstructure:"dest_addr"`
}

type Ethereum struct {
    RpcURL string `mapstructure:"rpc_url"`
    FromAccount EthereumAccount `mapstructure:"from_account"`
    ToAccount EthereumAccount `mapstructure:"to_account"`
}

type EthereumAccount struct {
    PrivateKey string `mapstructure:"private_key"`
    PublicKey string `mapstructure:"public_key"`
    Address string `mapstructure:"address"`
}
