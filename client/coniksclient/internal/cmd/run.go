package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/coniks-sys/coniks-go/client"
	"github.com/coniks-sys/coniks-go/eth"
	"github.com/coniks-sys/coniks-go/keyserver/testutil"
	p "github.com/coniks-sys/coniks-go/protocol"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

const help = "- register [name] [key]:\r\n" +
	"	Register a new name-to-key binding on the CONIKS-server.\r\n" +
	"- lookup [name]:\r\n" +
	"	Lookup the key of some known contact or your own bindings.\r\n" +
	"- enable timestamp:\r\n" +
	"	Print timestamp of format <15:04:05.999999999> along with the result.\r\n" +
	"- disable timestamp:\r\n" +
	"	Disable timestamp printing.\r\n" +
	"- help:\r\n" +
	"	Display this message.\r\n" +
	"- exit, q:\r\n" +
	"	Close the REPL and exit the client."

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run the test client.",
	Long:  "Run gives you a REPL, so that you can invoke commands to perform CONIKS operations including registration and key lookup. Currently, it supports:\n" + help,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd)
	},
}

func init() {
	RootCmd.AddCommand(runCmd)
	runCmd.Flags().StringP("config", "c", "config.toml",
		"Config file for the client (contains the server's initial public key etc).")
	runCmd.Flags().BoolP("debug", "d", false, "Turn on debugging mode")
	runCmd.Flags().StringP("ethconfig", "t", "eth.toml", "Path to ethereum configuration file")
	runCmd.Flags().BoolP("eth", "e", false, "Enable auditing with Ethereum")
}

func run(cmd *cobra.Command) {
	isDebugging, _ := strconv.ParseBool(cmd.Flag("debug").Value.String())
	conf := loadConfigOrExit(cmd)
	cc := p.NewCC(nil, true, conf.SigningPubKey)

	// Trusternity init for client

	auditEnabled, _ := strconv.ParseBool(cmd.Flag("eth").Value.String())
	ethConfig := cmd.Flag("ethconfig").Value.String()
	trustObject := eth.NewTrusternityObject(ethConfig)

	state, err := terminal.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	defer terminal.Restore(int(os.Stdin.Fd()), state)
	term := terminal.NewTerminal(os.Stdin, "coniks-client> ")
	for {
		line, err := term.ReadLine()
		if err != nil {
			writeLineInRawMode(term, err.Error(), isDebugging)
			return
		}

		args := strings.Fields(line)
		if len(args) < 1 {
			writeLineInRawMode(term, `[!] Type "help" for more information.`, isDebugging)
			continue
		}
		cmd := args[0]

		switch cmd {
		case "exit", "q":
			writeLineInRawMode(term, "[+] See ya.", isDebugging)
			return
		case "help":
			writeLineInRawMode(term, help, false) // turn off debugging mode for this command
		case "enable", "disable":
			if len(args) != 2 {
				writeLineInRawMode(term, "[!] Unrecognized command: "+line, isDebugging)
				continue
			}
			switch args[1] {
			case "timestamp":
				if cmd == "enable" {
					isDebugging = true
				} else {
					isDebugging = false
				}
			default:
				writeLineInRawMode(term, "[!] Unrecognized command: "+line, isDebugging)
			}
		case "register":
			if len(args) != 3 {
				writeLineInRawMode(term, "[!] Incorrect number of args to register.", isDebugging)
				continue
			}
			msg, _ := register(cc, conf, args[1], args[2])
			writeLineInRawMode(term, "[+] "+msg, isDebugging)
		case "lookup":
			if len(args) != 2 {
				writeLineInRawMode(term, "[!] Incorrect number of args to lookup.", isDebugging)
				continue
			}
			msg, _ := keyLookup(cc, conf, args[1])
			writeLineInRawMode(term, "[+] "+msg, isDebugging)
		case "audit":
			if auditEnabled {
				if len(args) != 2 {
					writeLineInRawMode(term, "[!] Incorrect number of args to lookup.", isDebugging)
					continue
				}
				epoch, _ := strconv.ParseUint(args[1], 10, 64)
				msg := trustObject.AuditSTR(epoch)
				writeLineInRawMode(term, "[+] Query epoch "+args[1]+". STR: "+msg, isDebugging)
			} else {
				writeLineInRawMode(term, "[+] Please enable ethereum mode", isDebugging)
			}
		default:
			writeLineInRawMode(term, "[!] Unrecognized command: "+cmd, isDebugging)
		}
	}
}

func createAuthRegistrationMsg(username string, accessToken string) ([]byte, error) {
	type Auth struct {
		Username    string
		AccessToken string
	}
	req := &Auth{Username: username, AccessToken: accessToken}
	return json.Marshal(req)
}

func register(cc *p.ConsistencyChecks, conf *client.Config, name string, key string) (string, p.ErrorCode) {
	var req []byte
	nameAccesstoken := strings.Split(name, " ")
	var accessToken string
	if len(nameAccesstoken) == 2 {
		name = nameAccesstoken[0]
		accessToken = nameAccesstoken[1]
		req, _ = createAuthRegistrationMsg(name, accessToken)
	} else {
		name = name
	}
	regMsg, err := client.CreateRegistrationMsg(name, []byte(key))
	if err != nil {
		return ("Couldn't marshal registration request!"), 500
	}
	if accessToken != "" {
		req = []byte(fmt.Sprintf(`{"Auth": %s, "ConiksRequest": %s}`, req, regMsg))
	} else {
		req = regMsg
	}

	var res []byte
	regAddress := conf.RegAddress
	if regAddress == "" {
		// fallback to conf.Address if empty
		regAddress = conf.Address
	}
	u, _ := url.Parse(regAddress)
	switch u.Scheme {
	case "https":
		res, err = testutil.NewHTTPSClient(req, regAddress)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	case "tcp":
		res, err = testutil.NewTCPClient(req, regAddress)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	case "unix":
		res, err = testutil.NewUnixClient(req, regAddress)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	default:
		return ("Invalid config!"), 500
	}

	response := client.UnmarshalResponse(p.RegistrationType, res)
	err = cc.HandleResponse(p.RegistrationType, response, name, []byte(key))
	switch err {
	case p.CheckBadSTR:
		// FIXME: remove me
		return ("Error: " + err.Error() + ". Maybe the client missed an epoch in between two commands, monitoring isn't supported yet."), p.CheckBadSTR
	case p.CheckPassed:
		switch response.Error {
		case p.ReqSuccess:
			return ("Succesfully registered name: " + name), p.ReqSuccess
		case p.ReqNameExisted:
			return ("Name is already registered."), p.ReqNameExisted
		}
	case p.CheckBindingsDiffer:
		switch response.Error {
		case p.ReqNameExisted:
			return (`Are you trying to update your binding? Unfortunately, KeyChange isn't supported yet.`), p.ReqNameExisted
		case p.ReqSuccess:
			recvKey, err := response.GetKey()
			if err != nil {
				return ("Oops! The server snuck in some other key. However, I cannot get the key from the response, error: " + err.Error()), 500
			}
			return ("Oops! The server snuck in some other key. [" + string(recvKey) + "] was registered instead of [" + string(key) + "]"), p.ReqSuccess
		}
	default:
		return ("Error: " + err.Error()), 500
	}
	return "", 500
}

func keyLookup(cc *p.ConsistencyChecks, conf *client.Config, name string) (string, p.ErrorCode) {
	req, err := client.CreateKeyLookupMsg(name)
	if err != nil {
		return ("Couldn't marshal key lookup request!"), 500
	}

	var res []byte
	u, _ := url.Parse(conf.Address)
	switch u.Scheme {
	case "https":
		res, err = testutil.NewHTTPSClient(req, conf.Address)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	case "tcp":
		res, err = testutil.NewTCPClient(req, conf.Address)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	case "unix":
		res, err = testutil.NewUnixClient(req, conf.Address)
		if err != nil {
			return ("Error while receiving response: " + err.Error()), 500
		}
	default:
		return ("Invalid config!"), 500
	}

	response := client.UnmarshalResponse(p.KeyLookupType, res)
	if key, ok := cc.Bindings[name]; ok {
		err = cc.HandleResponse(p.KeyLookupType, response, name, []byte(key))
	} else {
		err = cc.HandleResponse(p.KeyLookupType, response, name, nil)
	}
	switch err {
	case p.CheckBadSTR:
		// FIXME: remove me
		return ("Error: " + err.Error() + ". Maybe the client missed an epoch in between two commands, monitoring isn't supported yet."), p.CheckBadSTR
	case p.CheckPassed:
		switch response.Error {
		case p.ReqSuccess:
			key, err := response.GetKey()
			if err != nil {
				return ("Cannot get the key from the response, error: " + err.Error()), 500
			}
			return ("Found! Key bound to name is: [" + string(key) + "]"), p.ReqSuccess
		case p.ReqNameNotFound:
			return ("Name isn't registered."), p.ReqNameNotFound
		}
	default:
		return ("Error: " + err.Error()), 500
	}
	return "", 500
}
