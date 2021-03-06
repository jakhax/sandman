package sandbox;
import (

	"os"
	"io"
	"io/ioutil"
	"errors"
	"fmt"
	"regexp"
	// "strings"
	"encoding/json"
	"github.com/jakhax/sandman/runner"
	"github.com/jakhax/sandman/containers"
	"github.com/jakhax/sandman/runner/python"
	"github.com/jakhax/sandman/runner/golang"
	"github.com/jakhax/sandman/runneropt"
	"github.com/jakhax/sandman/spawn"
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

// MissingLanguageImage error 
type MissingLanguageImage struct{
	Language string
}

func (e MissingLanguageImage) Error() string{
	return fmt.Sprintf("Missing image for language: %s",e.Language)
}

// SandBox executes code in runner
type SandBox struct{

}

// Run method executes code in the sandbox
func (s *SandBox) Run(opt *runneropt.Opt) (stdout, stderr io.Reader, err error){
	image,err := GetRunnerImage(opt.Language);
	if err != nil{
		return;
	}
	containerService, err := containers.NewDockerSdkContainerService();
	
	if err !=  nil{
		return;
	}
	optJson,err := json.Marshal(opt)
	if err != nil{
		return
	}
	runContainerOption := containers.RunContainerOptions{
		Cmd:[]string{"run_json","-j",string(optJson)},
		Image:image,
		Runtime:"runsc-kvm",
		Timeout:opt.Timeout,
	}
	stdout, stderr, err = containerService.Run(runContainerOption);
	if err !=  nil{
		return;
	}
	
	// WriteToStd(stdout, stdout)
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

	image, ok :=  getImageFromConf(language); 
	if !ok{
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

	defer func(){
		WriteToStd(stdout,stderr)
		if err != nil{
			errMsg := fmt.Sprintf("Error: %s",err.Error())
			os.Stderr.WriteString(errMsg)
		}

	}()
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
	if len(opt.Shell) > 0{
		timeout,ok := opt.LanguagesConf.Timeouts["setup-shell"]
		if !ok{
			err = errors.New("Must Provide time out for setup shell");
			return
		}
		spawnOpt := & spawn.Opt{
			Dir:opt.Dir,
			Timeout:timeout,
		};
		shellStdout, shellStderr, errX := spawn.RunShell(spawnOpt,opt.Shell)
		if errX !=  nil{
			err = errX
			stdout = shellStdout
			stderr = shellStderr
			return
		}
	}
	
	// run strategy
	if opt.Strategy == SolutionOnlyStrategy {
		stdout,stderr,err = codeRunner.SolutionOnly(opt)
	}else{
		stdout,stderr, err = codeRunner.TestIntegration(opt)
	}
	if err != nil{
		return
	}
	//transform
	stdout,stderr, err = codeRunner.TransformOutput(stdout,stderr)
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
	return
}

type Test struct{
	Name string
	Passed *bool
	Failed *bool
	Error *bool
	Content string
}

func ParseTests(tests []string){

	for _,t := range tests{
		test := &Test{}
		//get test name
		pattern := regexp.MustCompile(`<IT::>(?P<Name>.*)`)
		m := pattern.FindStringSubmatch(t)
		pm := pattern.SubexpNames()
		if len(pm) > 0 && pm[0]=="Name"{
			test.Name = m[0]
		}

	}
}

func GetTests(text string)(tests []string){
	pattern := regexp.MustCompile(`<IT::>.*\n(?:.*\n)+?(.*?)<COMPLETEDIN::>[\d.]+`)
	matches := pattern.FindAllStringSubmatch(text,-1)
	if matches == nil{
		return
	}
	tests = []string{}
	for _,sm := range matches{
		tests = append(tests,sm[0])
	}
	return
}

//CleanTags removes tags from output since we dont use them
func CleanTags(text string)string{
	// mapper :=  make(map[string]string)
	// mapper[`<DESCRIBE::>`] = `Tests: `
	// mapper[`<IT::>`] = `Test: `
	// mapper[`<PASSED::>`] = ``
	// mapper[`<FAILED::>`] = ``
	// mapper[`<ERROR::>`] = `Error: `
	// mapper[`<COMPLETEDIN::>`] = `Completed In: `
	// for k,v := range mapper{
	// 	text := strings.ReplaceAll(text,k,v)
	// }
	// //remove log tag
	// pattern := regexp.MustCompile(`<LOG(\w|\:)+>`)
	// text = pattern.ReplaceAllString(text,"")
	// //remove more than 1 new line
	// pattern = regexp.MustCompile(`<:LF:>`)
	// text = pattern.ReplaceAllString(text,"\n")
	// pattern = regexp.MustCompile(`(\n){2,}`)
	// text = pattern.ReplaceAllString(text,"\n")
	// //remove more than 1 new line

	return text
}

// WriteToStd writes to stdin/stderr
func WriteToStd(stdout,stderr io.Reader) (err error){
	if stdout != nil{
		stdOutB ,errX := ioutil.ReadAll(stdout)
		if errX != nil {
			err = errX
			return
		}
		_, err = os.Stdout.WriteString(string(stdOutB))
		if err != nil{
			return
		}
	}
	if stderr != nil{
		stdErrB ,errX := ioutil.ReadAll(stderr)
		if errX != nil {
			err = errX
			return
		}
		_, err = os.Stderr.WriteString(string(stdErrB))
		return
	}
	return
}

// CreateCodeRunner returns a language specific code runner
func CreateCodeRunner(language string)(codeRunner runner.CodeRunner, err error){
	
	switch language {
		case runner.PYTHON:
			codeRunner = &python.Runner{};
			break;
		case runner.GOLANG:
			codeRunner = &golang.Runner{};
			break;
		default:
			err = errors.New("Missing Language Code Runner")
			break
	}
	return;
}

// SetupFromOpt setup for opt
func SetupFromOpt(opt *runneropt.Opt) (err error){
	// languages conf
	lc, err :=  runneropt.GetLanguagesConf();
	if err != nil {
		return
	}
	opt.LanguagesConf = lc 
	//working dir
	if opt.Dir == ""{
		wd,ok := opt.LanguagesConf.WorkingDir[opt.Language]
		if !ok{
			wd = opt.LanguagesConf.WorkingDir["default"]
		}
		opt.Dir = wd
	}
	// get strategy
	if len(opt.Fixture) > 0 {
		opt.Strategy = TestIntegrationStrategy
	}else{
		opt.Strategy = SolutionOnlyStrategy
	}
	//timeout 
	if opt.Timeout == 0{
		timeout,ok := opt.LanguagesConf.Timeouts[opt.Language]
		if !ok{
			timeout = opt.LanguagesConf.Timeouts["default"]
		}
		opt.Timeout = timeout
	}

	return;
}

// NewSandBoxRunner returns a sandbox runner
func NewSandBoxRunner() (sanboxRunner SandBoxRunner, err error){
	sanboxRunner = &SandBoxBaseRunner{}
	return
}