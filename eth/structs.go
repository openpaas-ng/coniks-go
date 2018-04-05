package eth

// A EtherConfig contains configuration values
// which are read at initialization time from
// a TOML format configuration file.
type EtherConfig struct {
	//Geth RPC endpoint url that we can use to connect to
	EndpointURL string `toml:"geth_rpc_endpoint_url"`
	//Ethereum Account address that is used to send transaction
	//This account should have enough gas to spent on transaction
	//format: hex string
	AccountAddress string `toml:"eth_account_address"`	
	//Ethereum contract address that we will send transaction to
	//format: hex string
	TrusternityContractAddress string `toml:"trusternity_contract_address"`
	configFilePath string
}

type transactionRPCRequest struct {
	ID      string
	Jsonrpc string
	Method  string
	Params  []transaction
}

type ethCallRPCRequest struct {
	ID      string
	Jsonrpc string
	Method  string
	Params  []interface{}
}

type ethCallParams struct {
	Trans transaction `json:"Object"`
	Param string `json:"param"`
}


type rpcResponse struct {
	Jsonrpc string
	Result  string
	ID      string
}

type transaction struct {
	From     string
	To       string
	Gas      string `json:",omitempty"`
	GasPrice string `json:",omitempty"`
	Value    string `json:",omitempty"`
	Data     string
}