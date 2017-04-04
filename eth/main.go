package main

import (
	"fmt"

	"encoding/hex"

	"github.com/coniks-sys/coniks-go/eth/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	testcrypto := crypto.Keccak256([]byte("abc"))
	fmt.Printf("%x\n", testcrypto)

	utils.RPCEndpointURL = "http://localhost:8545"
	utils.BaseAddress = "0x2589a613B33D98cC1f9b0874dB59fc315174DF58"
	utils.ContractAddress = "0x652f54e04b5961D983e88b36904F003A97A707A4"

	//res := utils.GetCoinBaseAddress("1")
	//res := utils.Register("1", "loria.fr", make([]byte, 4))

	//testPublish()	
	testGetSTR();
}

func testGetIPName(){
	res := utils.GetIPName("1", utils.BaseAddress)
	fmt.Println(res)
}

func testGetSTR(){
	res := utils.GetSTR("1", 1, utils.BaseAddress)
	fmt.Println(res)
}

func testPublish() {
	dummySTR, error := hex.DecodeString("17bce5eb73845f7dbb335ed15b2f39be22a36bd1e92c926107b5766ade5139a0")

	if error != nil {
		fmt.Println(error)
	}
	res := utils.Publish("2", 1, dummySTR)
	fmt.Println(res)
}
