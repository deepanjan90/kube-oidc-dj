package cmd

import (
    "io"
    "os"
	"log"
	"fmt"
    "bufio"
    "bytes"
    "os/exec"
    "strings"
    "path/filepath"

	"github.com/spf13/cobra"
)

type Config map[string]string

var getTokenCmd = &cobra.Command{
	Use:   "getToken",
	Short: "Get the kubernetes cluster token to connect.",
	Long: `This is a CLI library used to connect to a kubernetes cluster.
This application uses onelogin credentials to authenticate a user 
and generate an access token to connect to the cluster.`,
	Run: func(cmd *cobra.Command, args []string) {

		dirname, error := os.UserHomeDir()
	    if error != nil {
	        log.Fatal( error)
	    }

	    oneloginCredFile := filepath.Join(dirname, ".gyro", "onelogin", "credentials")

	    file, error := os.Open(oneloginCredFile)
	    if error != nil {
	        fmt.Println(error)
	    }
	    defer file.Close()

	    //
    	reader := bufio.NewReader(file)

    	config := Config{}

	    for {
	        line, error := reader.ReadString('\n')

	        // check if the line has = sign
	        // and process the line. Ignore the rest.
	        if equal := strings.Index(line, "="); equal >= 0 {
	            
	            if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
	                value := ""
	                if len(line) > equal {
	                    value = strings.TrimSpace(line[equal+1:])
	                }
	                // assign the config map
	                config[key] = value
	            }
	        }

	        if error == io.EOF {
	            break
	        }

	        if error != nil {
	            fmt.Println(error)
	            break
	        }
	    }

	    //


		lsCmd := exec.Command("kubectl",
				"oidc-login",
				"get-token",
				"--oidc-issuer-url=" + config["oidcIssuerUrl"],
				"--oidc-client-id=" + config["oidcClientId"],
				"--oidc-client-secret=" + config["oidcClientSecret"],
				"--oidc-extra-scope=" + config["oidcExtraScope"],
				"--grant-type=" + config["grantType"])

		var out bytes.Buffer
    	lsCmd.Stdout = &out

    	cmdError := lsCmd.Run()

    	if cmdError != nil {
        	log.Fatal(cmdError)
    	} else {
    		fmt.Println(out.String())
    	}
	},
}

func init() {
	rootCmd.AddCommand(getTokenCmd)
}
