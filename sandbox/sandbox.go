package sandbox;
import (

	"os"
	"io"
	"io/ioutil"
	"errors"
	"fmt"
	"github.com/jakhax/sandman/runner"
	"github.com/sirupsen/logrus"
	"github.com/jakhax/sandman/containers"
	"github.com/jakhax/sandman/runner/python"
	"github.com/jakhax/sandman/runneropt"
)

const (
	// SolutionOnlyStrategy strategy
	SolutionOnlyStrategy = "solutionOnly"
	// TestIntegrationStrategy strategy
	TestIntegrationStrategy = "testIntegration"
)

// SandBoxRunner interface
type SandBoxRunner interface{
	Run(opt *runneropt.Opt) (stdout,stderr io.Reader, err error);
}

// SandBoxInterface interface for sandbox
type SandBoxInterface interface{
	Run(opt *runneropt.Opt) (stdout, stderr io.Reader, err error)
}

// SandBox executes code in runner
type SandBox struct{

}

// MissingLanguageImage error 
type MissingLanguageImage struct{
	Language string
}

func (e MissingLanguageImage) Error() string{
	return fmt.Sprintf("Missing image for language: %s",e.Language)
}

// Run method executes code in the sandbox
func (s *SandBox) Run(opt *runneropt.Opt) (stdout, stderr io.Reader, err error){
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
	WriteToStd(stdOut, stdErr)
	return;
}

// NewSandBox returns a sandbox
func NewSandBox()(sandbox *SandBox, err error){
	sandbox  = &SandBox{}
	return;
}

//GetRunnerImage returns language runner image
func GetRunnerImage(language string) (image string, err error){
	languagesConf , err :=  runneropt.GetLanguagesConf()
	if err != nil{
		return
	}

	getImageFromConf :=  func (language string) (image string,ok bool){
		image, ok =  languagesConf.Images[language]
		return
	}

	switch language{
		case runner.PYTHON:
			image, _ =  getImageFromConf(language);
			break;
		default:
			break;
	}
	if image == ""{
		err =  MissingLanguageImage{Language:language}
	}
	return;
}

const (
	//HomePath default home directory in sandbox
	HomePath = "/home/appuser/"
)

//SandBoxBaseRunner for sandbox environment
type SandBoxBaseRunner struct{

}

//Run method
func (r *SandBoxBaseRunner) Run(opt *runneropt.Opt) (stdout,stderr io.Reader, err error) {
	// setup 
	err =  SetupFromOpt(opt)
	if err != nil{
		return
	}
	codeRunner, err := CreateCodeRunner(opt.Language)
	if err !=  nil {
		return
	}
	// run setup shell code if exists
	if opt.Shell != nil{
		timeout,ok := opt.LanguagesConf.Timeouts["setup-shell"]
		if !ok{
			err = errors.New("Must Provide time out for setup shell");
			return
		}
		spawnOpt := & SpawnOpt{
			Dir:opt.Dir,
			Timeout:timeout,
		};
		shellStdin, shellStderr, errX := RunShell(spawnOpt,opt.Shell)
		if errX !=  nil{
			err = errX
			WriteToStd(shellStdin,shellStderr)
			return
		}
	}
	// run strategy
	if opt.Strategy == SolutionOnlyStrategy {
		stdout,stdout, err = codeRunner.SolutionOnly(opt)
	}else{
		stdout,stdout, err = codeRunner.TestIntegration(opt)
	}
	if err != nil{
		return
	}
	//transform
	stdout,stdout, err = codeRunner.TransformOutput(stdout,stderr)
	if err != nil{
		return
	}
	// sanitize 
	stdout,err = codeRunner.SanitizeStdOut(stdout)
	if err != nil{
		return
	}
	// sanitize 
	stderr,err = codeRunner.SanitizeStdErr(stderr)
	if err != nil{
		return
	}
	WriteToStd(stdout,stdout)
	return
}

// WriteToStd writes to stdin/stderr
func WriteToStd(stdin,stderr io.Reader) (err error){
	stdInB ,err := ioutil.ReadAll(stdin)
	if err != nil {
		return
	}
	stdErrB ,err := ioutil.ReadAll(stderr)
	if err != nil {
		return
	}
	_, err = os.Stdout.WriteString(string(stdInB))
	if err != nil{
		return
	}
	_, err = os.Stderr.WriteString(string(stdErrB))
	return
}

// CreateCodeRunner returns a language specific code runner
func CreateCodeRunner(language string)(codeRunner runner.CodeRunner, err error){
	
	switch language {
		case runner.PYTHON:
			codeRunner = &python.Runner{};
			break;
		default:
			err = errors.New("Missing Language Code Runner")
			break
	}
	return;
}

// SetupFromOpt setup for opt
func SetupFromOpt(opt *runneropt.Opt) (err error){
	if opt.Dir == ""{
		wd,ok := opt.LanguagesConf.WorkingDir[opt.Language]
		if !ok{
			wd = opt.LanguagesConf.WorkingDir["default"]
		}
		opt.Dir = wd
	}
	// get strategy
	if opt.Fixture != nil {
		opt.Strategy = TestIntegrationStrategy
	}else{
		opt.Strategy = SolutionOnlyStrategy
	}
	// languages conf
	lc, err :=  runneropt.GetLanguagesConf();
	if err != nil {
		return
	}
	opt.LanguagesConf = lc 
	return;
}

