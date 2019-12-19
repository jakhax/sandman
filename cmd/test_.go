package cmd;
import (
	"io"
	// "os"
	"github.com/spf13/cobra"
	// "github.com/sirupsen/logrus"
	"github.com/jakhax/sandman/spawn"

)

var testCommand = &cobra.Command{
	Use:"test_",
	Short:"test ",
	RunE: func(cmd *cobra.Command, args []string) (err error){
		// stdout,stderr,err = codeRunner.SolutionOnly(opt)
		spawnOpt := &spawn.Opt{
			Dir:"/home/octojob",
			Timeout:200,
		};
		// code := string(opt.Code)
		var stdin io.Reader
		_, _, err = spawn.Spwan(spawnOpt,
			"python",[]string{"-c","while True:print(111)"},stdin)

		return;
	},
}



func init(){
	rootCmd.AddCommand(testCommand);
}