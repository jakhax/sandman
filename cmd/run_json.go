package cmd;
import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/jakhax/sandman/runners"
	"encoding/json"
)

var runJSONCommand = &cobra.Command{
	Use:"run_json",
	Short:"run from json input",
	RunE: func(cmd *cobra.Command, args []string) (err error){
		jsonInput,err :=  cmd.Flags().GetString("json");
		
		if err != nil {
			return
		}
		logrus.Info(jsonInput);
		opt := &runners.Opt{};

		err = json.Unmarshal([]byte(jsonInput),opt);
		if err != nil{
			return;
		}
		err =  opt.OK();
		return;
	},
}

func init(){
	runJSONCommand.Flags().StringP("json","j","","json input");
	runJSONCommand.MarkFlagRequired("json");
	rootCmd.AddCommand(runJSONCommand);
}