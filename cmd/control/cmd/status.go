package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/curusarn/resh/cmd/control/status"
	"github.com/curusarn/resh/pkg/msg"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "show RESH status",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("resh " + version)
		fmt.Println()
		fmt.Println("Resh versions ...")
		fmt.Println(" * installed: " + version + " (" + commit + ")")
		versionEnv, found := os.LookupEnv("__RESH_VERSION")
		if found == false {
			versionEnv = "UNKNOWN!"
		}
		commitEnv, found := os.LookupEnv("__RESH_REVISION")
		if found == false {
			commitEnv = "unknown"
		}
		fmt.Println(" * this shell session: " + versionEnv + " (" + commitEnv + ")")

		resp, err := getDaemonStatus(config.Port)
		if err != nil {
			fmt.Println(" * RESH-DAEMON IS NOT RUNNING")
			fmt.Println(" * Please REPORT this here: https://github.com/curusarn/resh/issues")
			fmt.Println(" * Please RESTART this terminal window")
			exitCode = status.Fail
			return
		}
		fmt.Println(" * daemon: " + resp.Version + " (" + resp.Commit + ")")

		if version != resp.Version || version != versionEnv {
			fmt.Println(" * THERE IS A MISMATCH BETWEEN VERSIONS!")
			fmt.Println(" * Please REPORT this here: https://github.com/curusarn/resh/issues")
			fmt.Println(" * Please RESTART this terminal window")
		}

		exitCode = status.ReshStatus
	},
}

func getDaemonStatus(port int) (msg.StatusResponse, error) {
	mess := msg.StatusResponse{}
	url := "http://localhost:" + strconv.Itoa(port) + "/status"
	resp, err := http.Get(url)
	if err != nil {
		return mess, err
	}
	defer resp.Body.Close()
	jsn, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error while reading 'daemon /status' response:", err)
	}
	err = json.Unmarshal(jsn, &mess)
	if err != nil {
		log.Fatal("Error while decoding 'daemon /status' response:", err)
	}
	return mess, nil
}
