package containers;

import (
	"io"
	"bytes"
	"fmt"
	"time"
	"os/exec"
	"github.com/jakhax/sandman/utils"
)

// DockerExecContainerService uses os installed docker to run container
type DockerExecContainerService struct{

}

// Run executes command in a docker container with a provided runtime
func (s *DockerExecContainerService) Run(runContainerOptions RunContainerOptions,
	) (stdOut io.Reader, stdErr io.Reader, err error){
		dockerCommand := []string{
			"run",
			"--rm",
		};

		// add runtime
		if(runContainerOptions.Runtime != ""){
			dockerCommand = append(dockerCommand, fmt.Sprintf("--runtime=%s",runContainerOptions.Runtime));
		}

		// volumes
		for k,v := range runContainerOptions.Volumes{
			dockerCommand = append(dockerCommand, fmt.Sprintf("-v %s:%s",k,v));
		}

		// network
		if runContainerOptions.DisableNetwork == false{
			containerNetwork := "none";
			dockerCommand = append(dockerCommand, fmt.Sprintf("--network=%s",containerNetwork));
		}
		
		// image
		dockerCommand = append(dockerCommand, runContainerOptions.Image);
		// command
		dockerCommand = append(dockerCommand, runContainerOptions.Cmd...);
		
		fmt.Println(dockerCommand);
		cmd := exec.Command("docker", dockerCommand...);
		logStdOut := new(bytes.Buffer);
		logStdErr := new(bytes.Buffer);
		cmd.Stdout = logStdOut;
		cmd.Stderr = logStdErr;

		// start cmd
		err = cmd.Start();
		if(err!=nil){
			return;
		}

		done :=  make(chan error);
		go func() { done <- cmd.Wait() }()
		// timer
		timeout := time.Duration(runContainerOptions.Timeout) * time.Second;
		timer := time.After(timeout);

		select{
			case err = <-done:
				if(err!=nil){
					err  =  &utils.ErrRunningContainer{Message:err.Error()};
					return;
				}
				stdOut = bytes.NewReader(logStdOut.Bytes());
				stdErr = bytes.NewReader(logStdErr.Bytes());
			case <-timer:
				err = cmd.Process.Kill();
				if(err!=nil){
					return;
				}
				err = &utils.ErrContainerTimeout{Message:"Error: Process Timed out"};
				return;
		}
		return;
	}

// NewDockerExecContainerService creates DockerExecContainerService{}
func NewDockerExecContainerService() (*DockerExecContainerService, error){
		s := &DockerExecContainerService{};
		return s,nil;
}
	
	