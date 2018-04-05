package eth

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

//GetCoinBaseAddress comment
func (trusternityObject *Trusternity) getCoinBaseAddress(id string) string {
	request := transactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_coinbase",
		ID:      id,
	}
	req, _ := json.Marshal(request)
	resp := CreatePostRequest(trusternityObject.config.EndpointURL, req)
	return resp
}

//Register function comment
func (trusternityObject *Trusternity) register(id, ipName string, fingerprint []byte) string {
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

	//make POST request
	request := makeEthSendRequest(trusternityObject.config, buffer.Bytes())
	resp := CreatePostRequest(trusternityObject.config.EndpointURL, request)

	return resp
}

//Publish function comment
func (trusternityObject *Trusternity) publish(id string, epoch uint64, STR []byte) string {
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

	transObject := transaction{
		From: trusternityObject.config.AccountAddress,
		To:   trusternityObject.config.TrusternityContractAddress,
		Data: transData,
		Gas:  gas,
	}

	request := transactionRPCRequest{
		Jsonrpc: "2.0",
		Method:  "eth_sendTransaction",
		ID:      id,
		Params:  []transaction{transObject},
	}
	req, _ := json.Marshal(request)
	log.Printf("Trusternity: publish %s", request)

	resp := CreatePostRequest(trusternityObject.config.EndpointURL, req)
	return resp
}

//GetIPName comment
func (trusternityObject *Trusternity) GetIPName(id, address string) string {
	// declare the main buffer
	var buffer bytes.Buffer
	// Init the buffer with method ID
	buffer.Write(getMethodID("GetProviderName(address)"))
	// its will be 4-32
	address = strings.Replace(address, "0x", "", 1)
	aBytes, _ := hex.DecodeString(address)
	buffer.Write(common.LeftPadBytes(aBytes, 32))

	//make POST request
	request := makeEthCallRequest(trusternityObject.config, buffer.Bytes())
	resp := CreatePostRequest(trusternityObject.config.EndpointURL, request)

	//parse the reply
	resp = strings.Replace(resp, "0x", "", 1)
	len := resp[64:len(resp)]
	name, _ := hex.DecodeString(len)

	return string(name)
}

// GetSTR comment
func (trusternityObject *Trusternity) GetSTR(id string, epoch uint64, address string) string {
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

	request := makeEthCallRequest(trusternityObject.config, buffer.Bytes())
	resp := CreatePostRequest(trusternityObject.config.EndpointURL, request)
	return resp
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
