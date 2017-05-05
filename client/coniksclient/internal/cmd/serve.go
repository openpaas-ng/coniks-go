package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"log"

	p "github.com/coniks-sys/coniks-go/protocol"
	"github.com/spf13/cobra"
)

// ServerCmd represents the base "testclient" command when called without any
// subcommands (register, lookup, ...).
var ServerCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the CONIKS client as a local http server",
	Long:  "Run the CONIKS client as a local http server, transferring the requests to the CONIKS Server",
	Run: func(cmd *cobra.Command, args []string) {
		startLocalHTTPServer(cmd)
	},
}

func init() {
	RootCmd.AddCommand(ServerCmd)
	ServerCmd.Flags().StringP("config", "c", "config.toml",
		"Config file for the client (contains the server's initial public key etc).")
}

func startLocalHTTPServer(cmd *cobra.Command) {
	http.HandleFunc("/", makeHandler(cmd))
	log.Println("Listening on http://localhost:3001")
	http.ListenAndServe(":3001", nil)
}

func makeHandler(cmd *cobra.Command) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)
		log.Printf("Received request : %s\n", bodyStr)
		conf := loadConfigOrExit(cmd)
		cc := p.NewCC(nil, true, conf.SigningPubKey)
		args := strings.Fields(bodyStr)
		if len(args) < 1 {
			msg := "[!] Usage :\n\nregister <name> <key>\nor\nlookup <name>\n"
			log.Print(msg)
			http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
			return
		}
		cmd := args[0]
		switch cmd {
		case "register":
			if len(args) != 3 {
				msg := "[!] Incorrect number of args to register.\nUsage : register <name> <key>\n"
				log.Printf(msg)
				http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
				return
			}
			msg, errCode := register(cc, conf, args[1], args[2])
			httpErrorCode := errorCodeToHTTPError(errCode)
			log.Printf("[+] Error code : %d, HTTP error code : %d, %s\n", errCode, httpErrorCode, msg)
			http.Error(w, fmt.Sprintf("[+] %s\n", msg), httpErrorCode)
		case "lookup":
			if len(args) != 2 {
				msg := "[!] Incorrect number of args to lookup.\nUsage : lookup <name>\n"
				log.Printf(msg)
				http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
				return
			}
			msg, errCode := keyLookup(cc, conf, args[1])
			httpErrorCode := errorCodeToHTTPError(errCode)
			log.Printf("[+] Error code : %d, HTTP error code : %d, %s\n", errCode, httpErrorCode, msg)
			http.Error(w, fmt.Sprintf("[+] %s\n", msg), httpErrorCode)
		default:
			log.Printf("[!] Unrecognized command: %s\n", cmd)
			http.Error(w, fmt.Sprintf("[!] Unrecognized command: %s\n", cmd), http.StatusBadRequest)
		}
	}
}

// Transform a CONIKS error code into HTTP Error code
// Success -> 200
// NameNotFound -> 404 Not Found
// NameAlreadyExists -> 409 Conflict
// CheckBadSTR -> 500 Internal Server error
// Other internal errors -> 500 Internal Server error
func errorCodeToHTTPError(errCode p.ErrorCode) int {
	var httpError int
	switch errCode {
	case p.ReqNameNotFound:
		httpError = http.StatusNotFound
	case p.ReqNameExisted:
		httpError = http.StatusConflict
	case p.ReqSuccess:
		httpError = http.StatusOK
	case p.CheckBadSTR, 500:
		httpError = http.StatusInternalServerError
	default:
		httpError = http.StatusInternalServerError
	}
	return httpError
}
