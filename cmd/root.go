package cmd;
import (
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
);

var rootCmd = &cobra.Command{
	Use:"Sandman",
	Short:"Root Command",
	Run:run,
}


func run(cmd *cobra.Command, args []string){
	cmd.Usage();
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}