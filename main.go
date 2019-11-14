package main
import (
	"fmt"
	"log"
	"io"
	"os"
	"strings"
	"io/ioutil"
	"github.com/jakhax/code-execution-engine/containers"
	"github.com/jakhax/code-execution-engine/utils"
	
)

func  play(){
	r:= strings.NewReader("Hello World");
	io.Copy(os.Stdout,r)
}

func main(){
	// play();
	containerService, err := containers.NewContainerService();
	utils.HandleError(err);
	runContainerOption := containers.RunContainerOptions{
		Cmd:[]string{"echo","Hello"},
		Image:"code-kombat/python-runner:latest",
		Runtime:"runsc-kvm",
	}
	stdOut, stdErr, err := containerService.RunContainer(runContainerOption);
	utils.HandleError(err);

	stdOut.Write

	io.Copy(os.Stdout, stdOut);
	io.Copy(os.Stderr, stdErr);

	var language Language;
	var source Source ;
	var strategy Strategy;
	language = "python3";
	source = "print('hello world')";
	strategy = "run";

	opt := Opt{
		Language : language,
		Source : source,
		Strategy: strategy,
	};

	var sandbox SandBoxInterface;
	sandbox = SandBox{
		Opt : opt,
	};

	output,err := ioutil.ReadAll(sandbox.Run());
	if(err != nil){
		log.Fatal(err);
	}
	
	results := string(output);

	fmt.Printf("%s\n",results);
	
}