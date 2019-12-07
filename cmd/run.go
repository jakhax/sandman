package cmd;
import (
	// "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/jakhax/sandman/runners"
	"github.com/jakhax/sandman/sandbox"
)

var runCommand = &cobra.Command{
	Use:"run",
	Short:"run code",
	RunE: func(cmd *cobra.Command, args []string) (err error){
		opt, err :=  CreateOptFromCommand(cmd);
		if err != nil{
			return
		}
		sanboxFlag , err := cmd.Flags().GetBool("sandbox");
		if err != nil{
			return;
		}
		if sanboxFlag{
			s,errX  :=  sandbox.NewSandBox();
			if err != nil{
				err = errX;
				return
			}
			_,_,err=s.Run(opt)
			if err != nil{
				// logrus.Error(err);
				return
			}
		}
		return;
	},
}

// CreateOptFromCommand returns opt from command
func CreateOptFromCommand(cmd *cobra.Command) (opt *runners.Opt, err error){
	opt = &runners.Opt{};
	code,err :=  cmd.Flags().GetString("code");
	if err != nil {
		return
	}
	opt.Code = code;
	language,err :=  cmd.Flags().GetString("language");
	if err != nil {
		return
	}
	opt.Language = language;
	setupCode,err :=  cmd.Flags().GetString("setup_code");
	if err != nil {
		return
	}
	opt.SetupCode = setupCode;
	fixture,err :=  cmd.Flags().GetString("fixture");
	if err != nil {
		return
	}
	opt.Fixture = fixture;
	testFramework,err :=  cmd.Flags().GetString("test_framework");
	if err != nil {
		return
	}
	opt.TestFramework = testFramework;
	format,err :=  cmd.Flags().GetString("output_format");
	if err != nil {
		return
	}
	opt.Format = format;
	timeout,err :=  cmd.Flags().GetInt("timeout");
	if err != nil {
		return
	}
	opt.Timeout = timeout;
	cpu,err :=  cmd.Flags().GetInt("cpu");
	if err != nil {
		return
	}
	opt.CPU = cpu;
	memory,err :=  cmd.Flags().GetInt("memory");
	if err != nil {
		return
	}
	opt.Memory = memory;
	shell,err :=  cmd.Flags().GetString("shell");
	if err != nil {
		return
	}
	opt.Shell = shell;
	// validate
	err = opt.OK()
	if err != nil{
		return;
	}
	return;
}

func init(){
	runCommand.Flags().Bool("sandbox",false,"environment to run code in, if set will execute code in sandbox");
	runCommand.Flags().StringP("code","c","","code to run");
	runCommand.Flags().StringP("setup_code","s","","Setup code to be used for executing the code");
	runCommand.Flags().StringP("fixture","f","","Test fixture code to test with");
	runCommand.Flags().StringP("test_framework","t","","Test framework to use");
	runCommand.Flags().StringP("language","l","","The language to execute the code in");
	runCommand.Flags().StringP("output_format","F","","Output format, options are 'default' and 'json'");
	runCommand.Flags().IntP("timeout","T",0,
	"The timeout to be used for running the code. If not specified a language specific default will be used");
	runCommand.Flags().IntP("cpu","C",0,"cpu limit")
	runCommand.Flags().IntP("memory","M",0,"memory limit")
	runCommand.Flags().StringP("shell","S","",
	"An optional shell script which will be ran within the sandbox environment before the code is executed");
	runCommand.MarkFlagRequired("code");
	runCommand.MarkFlagRequired("language");
	rootCmd.AddCommand(runCommand);
}