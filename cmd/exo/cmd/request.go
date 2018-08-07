package cmd

import (
	"fmt"

	"github.com/exoscale/egoscale"
)

// asyncRequest if no response expected send nil
func asyncRequest(cmd egoscale.AsyncCommand, msg string, response interface{}) error {
	fmt.Print(msg)
	var errorReq error
	cs.AsyncRequestWithContext(gContext, cmd, func(jobResult *egoscale.AsyncJobResult, err error) bool {

		fmt.Print(".")

		if err != nil {
			errorReq = err
			return false
		}

		if jobResult.JobStatus == egoscale.Success {
			if response != nil {
				if errR := jobResult.Result(response); errR != nil {
					errorReq = errR
				}
			}
			fmt.Println(" success.")
			return false
		}

		return true
	})

	if errorReq != nil {
		fmt.Println(" failure!")
	}

	return errorReq
}
