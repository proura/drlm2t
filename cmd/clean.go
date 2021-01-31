package cmd

import (
	"log"
	"os"

	"github.com/proura/drlm2t/model"
	"github.com/spf13/cobra"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("==CLEAN CALLED==")

		//Check for test to clean
		if len(args) < 1 {
			log.Fatalln("No test specified. Example: drlm2t clean example")
		} else {
			_, err := os.Stat("./tests/" + args[0])
			if os.IsNotExist(err) {
				log.Fatalln("Specified example \"", args[0], "\" does not exists")
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
