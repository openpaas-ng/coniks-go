package eth

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coniks-sys/coniks-go/crypto"
	"github.com/coniks-sys/coniks-go/protocol"
	"github.com/coniks-sys/coniks-go/utils"	
)

// Config stores globally for the ethereum package
var Config EtherConfig

// Initialize ethereum auditing package
func Initialize(configFile string) {
	Config, _ := loadConfig(configFile)
	_ = Config
}

// AuditSTR query STR of an epoch 
func AuditSTR(epoch uint64) string {		
	res := GetSTR("1", epoch, Config.AccountAddress)
	return res
}

// PublishSTR get the latest STR and publish to Ethereum
// This function should be called periodically after directory update at every Epoch
func PublishSTR(str *protocol.DirSTR) {
	digest := crypto.Digest(str.Signature)
	publish("1", str.Epoch, digest)
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

	conf.configFilePath = file
	return conf, nil
}
