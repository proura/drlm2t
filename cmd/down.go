package cmd

import (
	"log"
	"os"

	"github.com/proura/drlm2t/cfg"
	"github.com/proura/drlm2t/model"
	"github.com/spf13/cobra"
)

// downCmd represents the down command
var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop al VM and Networks",
	Long:  `Stop al VM and Networks`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("==DOWN CALLED==")

		//Check for test to down
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t down example")
		} else {
			_, err := os.Stat(cfg.Config.Drlm2tPath + "/tests/" + args[0])
			if os.IsNotExist(err) {
				log.Fatalln("Specified infrastructure", args[0], "does not exists")
			}
		}

		model.LoadRunningInfrastructure(args[0])
		model.Infrastructure = model.RunningInfrastructure

		model.Infrastructure.Status = "downing"
		model.SaveRunningIfrastructure()

		model.Infrastructure.DeleteHosts()
		model.Infrastructure.DeleteNetworks()

		model.Infrastructure.Status = "down"
		model.SaveRunningIfrastructure()
	},
}

func init() {
	rootCmd.AddCommand(downCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
