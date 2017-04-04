package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"io/ioutil"
	"net/http"
	"strings"
)

//RPCEndpointURL comment
var RPCEndpointURL string
var BaseAddress string
var ContractAddress string

//TransactionRPCRequest comment
type TransactionRPCRequest struct {
	ID      string
	Jsonrpc string
	Method  string
	Params  []Transaction
}

//EthCallRPCRequest comment
type EthCallRPCRequest struct {
	ID      string
	Jsonrpc string
	Method  string
	Params  []interface{}
}

//EthCallParams comment
type EthCallParams struct {
	Trans Transaction `json:"Object"`
	Param string `json:"param"`
}

//RPCResponse comment
type RPCResponse struct {
	Jsonrpc string
	Result  string
	ID      string
}

//Transaction comment
type Transaction struct {
	From     string
	To       string
	Gas      string `json:",omitempty"`
	GasPrice string `json:",omitempty"`
	Value    string `json:",omitempty"`
	Data     string
}

//CreatePostRequest create a http POST request to a host with a command body.
func CreatePostRequest(url string, command []byte) string {
	resp, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(command))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var rep RPCResponse
	json.Unmarshal(body, &rep)
	return rep.Result
}

//GetCoinBaseAddress comment
func GetCoinBaseAddress(id string) string {
	request := TransactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_coinbase",
		ID:      id,
	}
	req, _ := json.Marshal(request)
	resp := CreatePostRequest(RPCEndpointURL, req)
	return resp
}

//Register function comment
func Register(id, ipName string, fingerprint []byte) string {
	// declare the main buffer
	var buffer bytes.Buffer
	// Init the buffer with method ID
	buffer.Write(getMethodID("Register(string,bytes32)"))
	// its will be 4-32-32-data(32-32)-32 format
	// append location of the data part of the first paramenter
	loc1 := make([]byte, 4)
	binary.BigEndian.PutUint32(loc1, 64)
	buffer.Write(common.LeftPadBytes(loc1, 32))
	// append location of the data part of the second paramenter
	loc2 := make([]byte, 4)
	binary.BigEndian.PutUint32(loc2, 128)
	buffer.Write(common.LeftPadBytes(loc2, 32))
	// append the data part of the first argument
	// starts with the length of the bytes of the utf-8 encoded string
	domainBytes := []byte(ipName)
	length1 := make([]byte, 4)
	binary.BigEndian.PutUint32(length1, uint32(len(domainBytes)))
	buffer.Write(common.LeftPadBytes(length1, 32))
	// continue with the string
	buffer.Write(common.RightPadBytes(domainBytes, 32))
	// append the data part of the second argument
	// continue with the fingerprint
	buffer.Write(common.RightPadBytes(fingerprint, 32))
	//Put all data in hex
	transData := "0x" + hex.EncodeToString(buffer.Bytes())

	//Calculate the gas price - TODO: grab the estimate gas price from Geth
	gas := "0x" + fmt.Sprintf("%x", 900000)

	transObject := Transaction{	
		To:   ContractAddress,
		Data: transData,
		Gas:  gas,
	}

	request := TransactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendTransaction",
		ID:      id,
		Params:  []Transaction{transObject},
	}
	req, _ := json.Marshal(request)

	//s := string(req[:])
	//fmt.Println("[debug] ", s)

	resp := CreatePostRequest(RPCEndpointURL, req)
	return resp
}

func getMethodID(methodSignature string) []byte {
	digest := crypto.Keccak256([]byte(methodSignature))
	slice := digest[0:4] //get only first 4 bytes
	return slice
}

//Publish function comment
func Publish(id string, epoch uint64, STR []byte) string {
	// declare the main buffer
	var buffer bytes.Buffer
	// Init the buffer with method ID
	buffer.Write(getMethodID("Publish(uint64,bytes32)"))
	// its will be 4-32-32 format	
	// append the data part of the first argument
	// starts with the 64 bit uint epoch
	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, epoch)
	buffer.Write(common.LeftPadBytes(epochBytes, 32))

	// append the data part of the second argument
	// continue with the fingerprint
	buffer.Write(common.RightPadBytes(STR, 32))

	transData := "0x" + hex.EncodeToString(buffer.Bytes())

	//Calculate the gas price - TODO: grab the estimate gas price from Geth
	gas := "0x" + fmt.Sprintf("%x", 900000)

	transObject := Transaction{
		From: BaseAddress,
		To:   ContractAddress,
		Data: transData,
		Gas:  gas,
	}

	request := TransactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendTransaction",
		ID:      id,
		Params:  []Transaction{transObject},
	}
	req, _ := json.Marshal(request)

	//s := string(req[:])
	//fmt.Println("[debug] ", s)

	resp := CreatePostRequest(RPCEndpointURL, req)
	return resp
}

//GetIPName comment
func GetIPName(id, address string) string {
	// declare the main buffer
	var buffer bytes.Buffer
	// Init the buffer with method ID
	buffer.Write(getMethodID("GetProviderName(address)"))
	// its will be 4-32
	address = strings.Replace(address, "0x", "", 1)
	aBytes, _ := hex.DecodeString(address)
	buffer.Write(common.LeftPadBytes(aBytes, 32))

	transData := "0x" + hex.EncodeToString(buffer.Bytes())
	transObject := Transaction{
		From: BaseAddress,
		To:   ContractAddress,
		Data: transData,
	}

	request := EthCallRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		ID:      id,
		Params:   []interface{}{transObject, "latest"},
	}
	req, _ := json.Marshal(request)

	//s := string(req[:])
	//fmt.Println("[debug] ", s)

	resp := CreatePostRequest(RPCEndpointURL, req)
	resp = strings.Replace(resp, "0x", "", 1)

	len := resp[64:len(resp)]
	name,_ := hex.DecodeString(len)		

	return string(name)
}

// GetSTR comment
func GetSTR(id string, epoch uint64, address string) string{
	// declare the main buffer
	var buffer bytes.Buffer
	// Init the buffer with method ID
	buffer.Write(getMethodID("GetSTR(address,uint64)"))
	// its will be 4-32-32
	address = strings.Replace(address, "0x", "", 1)
	aBytes, _ := hex.DecodeString(address)
	buffer.Write(common.LeftPadBytes(aBytes, 32))

	epochBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(epochBytes, epoch)
	buffer.Write(common.LeftPadBytes(epochBytes, 32))

	transData := "0x" + hex.EncodeToString(buffer.Bytes())
	transObject := Transaction{
		From: BaseAddress,
		To:   ContractAddress,
		Data: transData,
	}

	request := EthCallRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_call",
		ID:      id,
		Params:   []interface{}{transObject, "latest"},
	}
	req, _ := json.Marshal(request)

	//s := string(req[:])
	//fmt.Println("[debug] ", s)

	resp := CreatePostRequest(RPCEndpointURL, req)	
	return resp
}