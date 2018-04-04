package eth

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func initTest() *Trusternity {
	conf := EtherConfig{
		AccountAddress:             "0xcaaC50c6D025A1F7F912f9d120f90796C6Fb30aA",
		EndpointURL:                "127.0.0.1:8545",
		TrusternityContractAddress: "0x0D9C76e58a12482E37c52023B9c34451B073D4F9",
	}

	var trusternityObject = new(Trusternity)
	trusternityObject.config = &conf
	return trusternityObject
}

// TestGetIPName comment
func TestGetIPName(t *testing.T) {
	trusternityObject := initTest()
	res := trusternityObject.GetIPName("1", trusternityObject.config.AccountAddress)
	fmt.Println(res)
}

// TestPublish comment
func TestPublish(t *testing.T) {
	trusternityObject := initTest()
	dummySTR, _ := hex.DecodeString("17bce5eb73845f7dbb335ed15b2f39be22a36bd1e92c926107b5766ade5139a0")
	res := trusternityObject.publish("9", 1, dummySTR)
	fmt.Println(res)
}

// TestGetSTR comment
func TestGetSTR(t *testing.T) {
	trusternityObject := initTest()
	res := trusternityObject.GetSTR("1", 1, trusternityObject.config.AccountAddress)
	fmt.Println(res)
}
