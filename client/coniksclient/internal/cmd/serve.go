package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
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

var allowedOrigins = []string{"127.0.0.1", "localhost"}

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
		if isOriginAllowed(r.Header.Get("Origin")) {
			log.Printf("Origin %s allowed\n", r.Header.Get("Origin"))
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}
		body, _ := ioutil.ReadAll(r.Body)
		bodyStr := string(body)
		log.Printf("Received request : %s\n", bodyStr)
		conf := loadConfigOrExit(cmd)
		cc := p.NewCC(nil, true, conf.SigningPubKey)
		args := strings.Fields(bodyStr)
		if len(args) < 1 {
			msg := "[!] Usage :\n\nregister <name> <key>\nor\nlookup <name>"
			log.Print(msg)
			http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
			return
		}
		cmd := args[0]
		switch cmd {
		case "register":
			if len(args) != 4 {
				msg := "[!] Incorrect number of args to register.\nUsage : register <name> <accessToken> <key>"
				log.Printf(msg)
				http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
				return
			}
			nameAccesstoken := fmt.Sprintf("%s %s", args[1], args[2])
			msg, errCode := register(cc, conf, nameAccesstoken, args[3])
			httpErrorCode := errorCodeToHTTPError(errCode)
			log.Printf("[+] Coniks protocol error code: %d - corresponding HTTP error code: %d - %s", errCode, httpErrorCode, msg)
			http.Error(w, fmt.Sprintf("[+] %s", msg), httpErrorCode)
		case "lookup":
			if len(args) != 2 {
				msg := "[!] Incorrect number of args to lookup.\nUsage : lookup <name>"
				log.Printf(msg)
				http.Error(w, fmt.Sprint(msg), http.StatusBadRequest)
				return
			}
			msg, errCode := keyLookup(cc, conf, args[1])
			httpErrorCode := errorCodeToHTTPError(errCode)
			log.Printf("[+] Coniks protocol error code: %d - corresponding HTTP error code: %d - %s", errCode, httpErrorCode, msg)
			http.Error(w, fmt.Sprintf("[+] %s", msg), httpErrorCode)
		default:
			log.Printf("[!] Unrecognized command: %s", cmd)
			http.Error(w, fmt.Sprintf("[!] Unrecognized command: %s", cmd), http.StatusBadRequest)
		}
	}
}

// Transform a CONIKS protocol error code into HTTP Error code
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

func isOriginAllowed(origin string) bool {
	allowedOriginsJoined := strings.Join(allowedOrigins, "|")
	var pattern = regexp.MustCompile(fmt.Sprintf(`(https?:\/\/)(%s)(:)([0-9]+)`, allowedOriginsJoined))

	return pattern.MatchString(origin)
}
