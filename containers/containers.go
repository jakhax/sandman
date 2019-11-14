package containers;

import (
	"io"
	"context"
	// "bytes"
	"fmt"
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
	Timeout *int
	Runtime string
}

type ContainerServiceInterface interface{
	BuildDockerImage(dockerFile string, containerName string) (err error);
	RunContainer(runContainerOptions RunContainerOptions) (stdOut io.Writer, stdErr io.Writer, err error)
}

type ContainerService struct{
	Client *client.Client
}


func (s *ContainerService) RunContainer(runContainerOptions RunContainerOptions,
	)(stdOut io.Writer, stdErr io.Writer, err error){

	ctx := context.Background();
	containerConfig := &container.Config{
		Image:runContainerOptions.Image,
		Cmd:runContainerOptions.Cmd,
		Volumes:runContainerOptions.Volumes,

	}
	if(runContainerOptions.Timeout != nil){
		containerConfig.StopTimeout = runContainerOptions.Timeout;
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
	
	err = s.Client.ContainerStart(
		ctx, 
		cnt.ID, 
		types.ContainerStartOptions{},
	);
	if(err != nil){
		return stdOut, stdErr, err;
	}

	statusCh, errCh := s.Client.ContainerWait(
		ctx,
		cnt.ID,
		container.WaitConditionNotRunning,
	)
	select{
		case err = <-errCh:
			if(err != nil){
				return stdOut, stdErr, err;
			}
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
		return stdOut, stdErr, err;
	}
	stdcopy.StdCopy(stdOut, stdErr, logs);
	return stdOut, stdErr, nil;

}

func (s ContainerService) BuildDockerImage(dockerFile string, 
	containerName string) (err error){
		return fmt.Errorf("Not Implemented");
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