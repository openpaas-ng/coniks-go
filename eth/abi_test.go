package eth

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func initializeTest() {
	Config = EtherConfig{
		AccountAddress:             "0x2589a613B33D98cC1f9b0874dB59fc315174DF58",
		EndpointURL:                "http://localhost:8545",
		TrusternityContractAddress: "0x652f54e04b5961D983e88b36904F003A97A707A4",
	}
}

// TestGetIPName comment
func TestGetIPName(t *testing.T) {
	initializeTest();
	res := GetIPName("1", Config.AccountAddress)
	fmt.Println(res)
}

// TestGetSTR comment
func TestGetSTR(t *testing.T) {
	initializeTest();
	res := GetSTR("1", 1, Config.AccountAddress)
	fmt.Println(res)
}

// TestPublish comment
func TestPublish(t *testing.T) {
	initializeTest();
	dummySTR, _ := hex.DecodeString("17bce5eb73845f7dbb335ed15b2f39be22a36bd1e92c926107b5766ade5139a0")
	res := publish("2", 1, dummySTR)
	fmt.Println(res)
}
