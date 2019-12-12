package cmd;
import (
	"github.com/spf13/cobra"
	"github.com/jakhax/sandman/runneropt"
	"github.com/jakhax/sandman/sandbox"
	"encoding/json"
)

var runJSONCommand = &cobra.Command{
	Use:"run_json",
	Short:"run from json input inside container",
	Run: func(cmd *cobra.Command, args []string){
		jsonInput,err :=  cmd.Flags().GetString("json");
		
		if err != nil {
			return
		}
		opt := &runneropt.Opt{};

		err = json.Unmarshal([]byte(jsonInput),opt);
		if err != nil{
			return;
		}
		err =  opt.OK();
		if err != nil{
			return
		}
		s,err :=  sandbox.NewSandBoxRunner();
		if err != nil{
			return
		}
		s.Run(opt)		
	},
}

func init(){
	runJSONCommand.Flags().StringP("json","j","","json input");
	runJSONCommand.MarkFlagRequired("json");
	rootCmd.AddCommand(runJSONCommand);
}