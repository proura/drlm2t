package cmd

import (
	"context"
	"log"
	"sync"
	"time"

	testserver "github.com/proura/drlm2t/drlm2t-server"
	"github.com/proura/drlm2t/model"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run all tests",
	Long:  `Run all tests`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t run example")
		}

		model.LoadRunningInfrastructure(args[0])
		model.Infrastructure = model.RunningInfrastructure

		// Arranquem el servidor per enviar i rebre les configs
		httpServerExitDone := &sync.WaitGroup{}
		httpServerExitDone.Add(1)
		srv := testserver.RunServer("testing", httpServerExitDone)

		if model.Infrastructure == nil {
			log.Fatalln("- There is nothing to run!")
		}

		for !model.Infrastructure.AllDone() {
			time.Sleep(1 * time.Second)
		}

		log.Println("+ ALL TESTS DONE.")

		if err := srv.Shutdown(context.TODO()); err != nil {
			log.Println(err.Error())
		}
		httpServerExitDone.Wait()
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
