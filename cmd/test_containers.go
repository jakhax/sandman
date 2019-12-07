package cmd;
import (
	"io"
	"os"
	"github.com/spf13/cobra"
	"github.com/sirupsen/logrus"
	"github.com/jakhax/sandman/containers"

)

var testContainerCommand = &cobra.Command{
	Use:"test_container",
	Short:"test container",
	RunE: func(cmd *cobra.Command, args []string) (err error){
		tesContainer()
		return;
	},
}


func tesContainer(){
	containerService, err := containers.NewDockerExecContainerService();
	
	if err !=  nil{
		logrus.Error(err);
		return;
	}
	runContainerOption := containers.RunContainerOptions{
		Cmd:[]string{"python","-c","while True:print(1)"},
		Image:"python:3.7.1-alpine3.7",
		Runtime:"runsc-kvm",
		Timeout:5,
	}
	stdOut, stdErr, err := containerService.Run(runContainerOption);
	if err !=  nil{
		logrus.Error(err);
		return;
	}
	io.Copy(os.Stdout, stdOut);
	io.Copy(os.Stderr, stdErr);
	if err !=  nil{
		logrus.Error(err);
		return;
	}
	
}

func init(){
	rootCmd.AddCommand(testContainerCommand);
}