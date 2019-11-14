// package main;

// import(
// 	"os"
// 	// "io"
// 	"context"
// 	"log"
// 	"runtime"
// 	// "fmt"
// 	"github.com/docker/docker/api/types"
// 	"github.com/docker/docker/api/types/container"
// 	"github.com/docker/docker/pkg/stdcopy"
// 	"github.com/docker/docker/client"
	
// )


// func handleError(err error){
// 	if(err != nil){
// 		_, fn, line, _ := runtime.Caller(1)
//         log.Fatalf("[error] %s:%d %v", fn, line, err)
// 	}
// }


// func RunContainer(image string){

	

// 	ctx := context.Background();

// 	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation());
// 	handleError(err);

// 	containerConfig := &container.Config{
// 		Image:image,
// 		Cmd:[]string{"echo","'Hello'"},
// 		Volumes
// 	};

// 	cont , err := cli.ContainerCreate(ctx,containerConfig,
// 		&container.HostConfig{},nil,"");
// 	handleError(err);

// 	err = cli.ContainerStart(ctx, cont.ID,types.ContainerStartOptions{});
// 	handleError(err);
// 	statusCh, errCh := cli.ContainerWait(ctx, cont.ID,container.WaitConditionNotRunning);
// 	select{
// 	case err:=<-errCh:
// 		handleError(err);
// 	case <-statusCh:
// 	}
// 	containerLogOptions := types.ContainerLogsOptions{
// 		ShowStderr:true,
// 		ShowStdout:true,
// 	}
// 	logs, err := cli.ContainerLogs(ctx, cont.ID, containerLogOptions);
// 	handleError(err);
// 	_, err = stdcopy.StdCopy(os.Stdout, os.Stderr, logs);
// 	handleError(err);
// }