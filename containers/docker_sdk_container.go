package containers;

import (
	"io"
	"context"
	"bytes"
	"time"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/docker/client"
	"github.com/jakhax/sandman/utils"
)


// DockerSdkContainerService uses go docker sdk
type DockerSdkContainerService struct{
	Client *client.Client
}

// Run executes command in a docker container with a provided runtime
func (s *DockerSdkContainerService) Run(runContainerOptions RunContainerOptions,
	)(stdOut io.Reader, stdErr io.Reader, err error){

	ctx := context.Background();
	containerConfig := &container.Config{
		Image:runContainerOptions.Image,
		Cmd:runContainerOptions.Cmd,
		Volumes:runContainerOptions.Volumes,
		NetworkDisabled:runContainerOptions.DisableNetwork,
		StopSignal:"SIGKILL",
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
		return;
	}
	err = s.Client.ContainerStart(
		ctx, 
		cnt.ID, 
		types.ContainerStartOptions{},
	);
	
	if(err != nil){
		return;
	}

	// timer
	timeout := time.Duration(runContainerOptions.Timeout) * time.Second;
	timer := time.After(timeout);

	statusCh, errCh := s.Client.ContainerWait(
		ctx,
		cnt.ID,
		container.WaitConditionNextExit,
	)
	
	select{
		case err = <-errCh:
			if(err != nil){
				err = &utils.ErrRunningContainer{Message:err.Error()};
				return;
			}
		case <- timer:
			stopTimeout := time.Duration(1) * time.Millisecond;
			err = s.Client.ContainerStop(ctx, cnt.ID, &stopTimeout);
			if err != nil {
				err = &utils.ErrRunningContainer{Message:err.Error()};
				return;
			}
			err = &utils.ErrContainerTimeout{Message:"Error: Process Timed out"}
			return;
		case <- statusCh:
	}

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
		return;
	}
	logStdOut := new(bytes.Buffer);
	logStdErr := new(bytes.Buffer);
	stdcopy.StdCopy(logStdOut, logStdErr, logs);
	stdOut = bytes.NewReader(logStdOut.Bytes())
	stdErr = bytes.NewReader(logStdErr.Bytes())
	return;
}

// NewDockerSdkContainerService creates DockerSdkContainerService{}
func NewDockerSdkContainerService() (*DockerSdkContainerService, error){
		s := &DockerSdkContainerService{};
		
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation());
		if(err != nil){
			return s,err;
		}
		s.Client = cli;
		return s,err;
	}

