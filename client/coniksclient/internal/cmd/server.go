package cmd

import (	
	"github.com/spf13/cobra"
	"net/http"
	"io"
)

func init() {
	RootCmd.AddCommand(ServerCmd)
}

// ServerCmd represents the base "testclient" command when called without any
// subcommands (register, lookup, ...).
var ServerCmd = &cobra.Command{
	Use:   "server",
	Short: "CONIKS client as local http server",
	Long:  "CONIKS client as local http server, it accepts only request from the browser",
	Run: func(cmd *cobra.Command, args []string) {
		startLocalHTTPServer()
	},
}

func startLocalHTTPServer() {	
	http.HandleFunc("/", handler)
    http.ListenAndServe(":3001", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}
