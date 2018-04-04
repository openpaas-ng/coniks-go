package eth

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
)

//CreatePostRequest create a http POST request to a host with a command body.
func CreatePostRequest(url string, command []byte) string {
	resp, err := http.Post("http://"+url, "application/json; charset=utf-8", bytes.NewBuffer(command))
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var rep rpcResponse
	json.Unmarshal(body, &rep)
	log.Printf("Trusternity: %s", rep.Result)
	return rep.Result
}

func getMethodID(methodSignature string) []byte {
	digest := crypto.Keccak256([]byte(methodSignature))
	slice := digest[0:4] //get only first 4 bytes
	return slice
}

func makeEthSendRequest(config *EtherConfig, data []byte) []byte {
	//Put all data in hex
	transData := "0x" + hex.EncodeToString(data)

	//Calculate the gas price - TODO: grab the estimate gas price from Geth
	gas := "0x" + fmt.Sprintf("%x", 900000)

	transObject := transaction{
		From: config.AccountAddress,
		To:   config.TrusternityContractAddress,
		Data: transData,
		Gas:  gas,
	}

	request := transactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendTransaction",
		ID:      "1",
		Params:  []transaction{transObject},
	}
	req, _ := json.Marshal(request)
	log.Printf("Trusternity: makeEthSendRequest %s", request)
	return req
}

func makeEthCallRequest(config *EtherConfig, data []byte) []byte {
	transData := "0x" + hex.EncodeToString(data)
	transObject := transaction{
		From: config.AccountAddress,
		To:   config.TrusternityContractAddress,
		Data: transData,
	}

	request := ethCallRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		ID:      "1",
		Params:  []interface{}{transObject, "latest"},
	}
	req, _ := json.Marshal(request)
	log.Printf("Trusternity: makeEthCallRequest %s", request)

	return req
}
