package sandbox;
import (
	"os"
	"io"
	// "io/ioutil"
	// "encoding/json"
	"github.com/jakhax/sandman/runners"
	"github.com/jakhax/sandman/utils"
	"github.com/sirupsen/logrus"
	"github.com/jakhax/sandman/containers"
)

// SandBoxInterface interface for sandbox
type SandBoxInterface interface{
	Run(opt *runners.Opt) (stdout, stderr io.Reader, err error)
}

// SandBox executes code in runner
type SandBox struct{

}

// Run method executes code in the sandbox
func (s *SandBox) Run(opt *runners.Opt) (stdout, stderr io.Reader, err error){
	if err != nil{
		return 
	}
	image,err := GetRunnerImage(opt.Language);
	if err != nil{
		return;
	}
	containerService, err := containers.NewDockerSdkContainerService();
	
	if err !=  nil{
		logrus.Error(err);
		return;
	}
	runContainerOption := containers.RunContainerOptions{
		Cmd:[]string{"python","-c","while True:print(1)"},
		Image:image,
		Runtime:"runsc-kvm",
		Timeout:opt.Timeout,
	}
	stdOut, stdErr, err := containerService.Run(runContainerOption);
	if err !=  nil{
		logrus.Error(err);
		return;
	}
	io.Copy(os.Stdout, stdOut);
	io.Copy(os.Stderr, stdErr);
	// stdout = output.StdOut;
	// stderr = output.StdErr;
	return;
}

// NewSandBox returns a sandbox
func NewSandBox()(sandbox *SandBox, err error){
	sandbox  = &SandBox{}
	return;
}



//GetRunnerImage returns language runner image
func GetRunnerImage(language string) (image string, err error){
	switch language{
		case "python":
			image = "sandman/python-runner";
			break;
		default:
			break;
	}
	if image == ""{
		err = utils.ValidationError{Message:"Language Runner image not found"}
	}
	return;
}