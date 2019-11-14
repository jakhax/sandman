package containers;

import (
	"io"
	"context"
	"bytes"
	"fmt"
	"time"
	"os/exec"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/client"
)

type RunContainerOptions struct{
	Name string
	Cmd []string
	Image string
	Volumes map[string]struct{}
	Timeout int
	Runtime string
}

type ContainerServiceInterface interface{
	GetClient() client.Client;
	BuildDockerImage(dockerFile string, containerName string) (err error);
	RunContainer(runContainerOptions RunContainerOptions) (stdOut io.Reader, stdErr io.Reader, err error)
}

type ExecContainerServiceInterface interface{
	RunContainer(runContainerOptions RunContainerOptions) (stdOut io.Reader, stdErr io.Reader, err error)
}

type ContainerService struct{
	Client *client.Client
}

type ExecContainerService struct{

}


func (s *ContainerService) GetClient() client.Client{
	return *s.Client;
}


func (s *ContainerService) RunContainer(runContainerOptions RunContainerOptions,
	)(stdOut io.Reader, stdErr io.Reader, err error){

	ctx := context.Background();
	containerConfig := &container.Config{
		Image:runContainerOptions.Image,
		Cmd:runContainerOptions.Cmd,
		Volumes:runContainerOptions.Volumes,
		NetworkDisabled:true,
		StopSignal:"SIGTERM",

	}
	if(&runContainerOptions.Timeout != nil){
		containerConfig.StopTimeout = &runContainerOptions.Timeout;
	}
	cnt, err := s.Client.ContainerCreate(
		ctx,
		containerConfig,
		&container.HostConfig{
			Runtime:runContainerOptions.Runtime,
		},
		nil,
		runContainerOptions.Name,
	);
	if(err != nil){
		return stdOut, stdErr, err;
	}
	fmt.Println("reached here");
	err = s.Client.ContainerStart(
		ctx, 
		cnt.ID, 
		types.ContainerStartOptions{},
	);
	
	if(err != nil){
		return stdOut, stdErr, err;
	}

	// statusCh, errCh := s.Client.ContainerWait(
	// 	ctx,
	// 	cnt.ID,
	// 	container.WaitConditionNextExit,
	// )
	// select{
	// 	case err = <-errCh:
	// 		if(err != nil){
	// 			return stdOut, stdErr, err;
	// 		}
	// 	case <- statusCh:
	// }

	timeout := time.Second * 5;

	
	err =  s.Client.ContainerStop(ctx, cnt.ID, &timeout);
	if(err != nil){
		return stdOut, stdErr, err;
	}
	// fmt.Println("reached here");

	logOptions := types.ContainerLogsOptions{
		ShowStderr:true,
		ShowStdout:true,
	};

	logs, err := s.Client.ContainerLogs(
		ctx,
		cnt.ID,
		logOptions,
	);
	if(err != nil){
		return stdOut, stdErr, err;
	}
	logStdOut := new(bytes.Buffer);
	logStdErr := new(bytes.Buffer);
	stdcopy.StdCopy(logStdOut, logStdErr, logs);
	stdOut = bytes.NewReader(logStdOut.Bytes())
	stdErr = bytes.NewReader(logStdErr.Bytes())
	return stdOut, stdErr, nil;

}

func (s ContainerService) BuildDockerImage(dockerFile string, 
	containerName string) (err error){
		return fmt.Errorf("Not Implemented");
	}


func (s *ExecContainerService) RunContainer(runContainerOptions RunContainerOptions,
	) (stdOut io.Reader, stdErr io.Reader, err error){

		fmt.Sprintf("docker run --rm --runtime")

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
			return stdOut, stdErr, err;
		}

		done :=  make(chan error);
		go func() { done <- cmd.Wait() }()
		// timer
		timeout := time.Duration(runContainerOptions.Timeout) * time.Second;
		timer := time.After(timeout);

		select{
			case err = <-done:
				stdOut = bytes.NewReader(logStdOut.Bytes());
				stdErr = bytes.NewReader(logStdErr.Bytes());
				if(err!=nil){
					return stdOut, stdErr, err;
				}
				// stdOut = bytes.NewReader(logStdOut.Bytes());
				// stdErr = bytes.NewReader(logStdErr.Bytes());

			case <-timer:
				err = cmd.Process.Kill();
				if(err!=nil){
					return stdOut, stdErr, err;
				}
				stdOut = bytes.NewReader([]byte(""));
				stdErr =  bytes.NewReader([]byte("\nError: Process Timed out\n"));
		}
		return stdOut, stdErr, err;
	}

func NewContainerService() (*ContainerService, error){
		s := &ContainerService{};
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation());
		if(err != nil){
			return s,err;
		}
		s.Client = cli;
		return s,err;
	}

func NewExecContainerService() (*ExecContainerService, error){
		s := &ExecContainerService{};
		return s,nil;
	}
	
	