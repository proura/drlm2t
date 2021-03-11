package cmd

import (
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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t run example")
		}

		model.LoadRunningInfrastructure(args[0])
		model.Infrastructure = model.RunningInfrastructure

		wg := new(sync.WaitGroup)
		wg.Add(1)

		go func() {
			testserver.RunServer("testing")
			wg.Done()
		}()

		if model.Infrastructure == nil {
			log.Fatalln("- There is nothing to run!")
		}

		for !model.Infrastructure.AllDone() {
			time.Sleep(1 * time.Second)
		}

		log.Println("+ ALL TESTS DONE. (ctrl+C to exit)")

		wg.Wait()

	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
