package eth

import (
	"log"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/coniks-sys/coniks-go/crypto"
	"github.com/coniks-sys/coniks-go/protocol"
	"github.com/coniks-sys/coniks-go/utils"
)

// A Trusternity object contains all auditing logics using Ethereum.
// When running CONIKS server with eth parameter, we create a new Trusternity
// object with configurations regarding the local Ethereum wallet and the RPC
// server.
type Trusternity struct {
	config *EtherConfig
}

// NewTrusternityObject Initialize an Trusternity object
func NewTrusternityObject(configFile string) *Trusternity {
	obj := new(Trusternity)
	conf, _ := loadConfig(configFile)
	obj.config = conf
	return obj
}

// AuditSTR query STR from the corresponding account address of the Trusternity object
func (trusternityObject *Trusternity) AuditSTR(epoch uint64) string {
	address := trusternityObject.config.AccountAddress
	res := trusternityObject.GetSTR("1", epoch, address)
	return res
}

// PublishSTR get the latest STR and publish to Ethereum
// This function should be called periodically after directory update at every Epoch
func (trusternityObject *Trusternity) PublishSTR(str *protocol.DirSTR) {	
	digest := crypto.Digest(str.Signature)
	trusternityObject.publish("1", str.Epoch, digest)
	log.Printf("Trusternity: Publish at epoch %d Digest %s", str.Epoch, digest)
}

// LoadConfig loads Ethereum configuration from eth.toml
func loadConfig(file string) (*EtherConfig, error) {
	var conf *EtherConfig
	if _, err := toml.DecodeFile(file, &conf); err != nil {
		return nil, fmt.Errorf("Failed to load Ethereum config: %v", err)
	}

	// load signing key
	conf.AccountAddress = utils.ResolvePath(conf.AccountAddress, file)
	conf.EndpointURL = utils.ResolvePath(conf.EndpointURL, file)
	conf.TrusternityContractAddress = utils.ResolvePath(conf.TrusternityContractAddress, file)

	log.Printf("Trusternity: Sucessfully loaded EthConfig")
	log.Printf("Trusternity: Account: %s", conf.AccountAddress)
	log.Printf("Trusternity: Contract: %s", conf.TrusternityContractAddress)

	conf.configFilePath = file
	return conf, nil
}
