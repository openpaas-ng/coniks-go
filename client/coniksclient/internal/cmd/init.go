package cmd

import (
	"fmt"
	"path"
	"strconv"

	"bytes"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/coniks-sys/coniks-go/client"
	"github.com/coniks-sys/coniks-go/utils"
	"github.com/spf13/cobra"
	//"github.com/tmp/keyserver/testutil"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a config file for the client.",
	Long: `Creates a file config.toml in the current working directory with
the following content:

sign_pubkey_path = "../../keyserver/coniksserver/sign.pub"
registration_address = "tcp://127.0.0.1:3000"
address = "tcp://127.0.0.1:3000"

[server-address]
address = "https://localhost:3001"
cert = "server.pem"
key = "server.key"

If the keyserver's public keys are somewhere else, you will have to modify the
config file accordingly.
`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := cmd.Flag("dir").Value.String()
		mkConfigOrExit(dir)
		cert, err := strconv.ParseBool(cmd.Flag("cert").Value.String())
		if err == nil && cert {
			//testutil.CreateTLSCert(dir)
		}
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("dir", "d", ".",
		"Location of directory for storing generated files")
	initCmd.Flags().BoolP("cert", "c", false, "Generate self-signed ssl keys/cert with sane defaults")
}

func mkConfigOrExit(dir string) {
	file := path.Join(dir, "config.toml")
	var conf = client.Config{
		SignPubkeyPath: "../../keyserver/coniksserver/sign.pub",
		RegAddress:     "tcp://127.0.0.1:3000",
		Address:        "tcp://127.0.0.1:3000",
	}

	var confBuf bytes.Buffer
	enc := toml.NewEncoder(&confBuf)
	if err := enc.Encode(conf); err != nil {
		fmt.Println("Coulnd't encode config. Error message: [" +
			err.Error() + "]")
		os.Exit(-1)
	}
	if err := utils.WriteFile(file, confBuf.Bytes(), 0644); err != nil {
		fmt.Println("Coulnd't write config. Error message: [" +
			err.Error() + "]")
		os.Exit(-1)
	}
}
