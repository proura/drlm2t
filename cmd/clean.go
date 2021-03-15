package cmd

import (
	"log"
	"os"

	"github.com/proura/drlm2t/cfg"
	"github.com/proura/drlm2t/model"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove al vm and test results",
	Long:  `Remove al vm and test results`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("==CLEAN CALLED==")

		//Check for test to clean
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t clean example")
		} else {
			_, err := os.Stat(cfg.Config.Drlm2tPath + "/tests/" + args[0])
			if os.IsNotExist(err) {
				log.Fatalln("Specified infrastructure \"", args[0], "\" does not exists")
			}
		}

		model.LoadRunningInfrastructure(args[0])
		model.Infrastructure = model.RunningInfrastructure
		model.Infrastructure.Clean()
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cleanCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cleanCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
