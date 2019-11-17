package containers;

import (
	"io"
)
// RunContainerOptions is the config needed to run code in a container
type RunContainerOptions struct{
	Name string
	Cmd []string
	Image string
	Volumes map[string]struct{}
	Timeout int
	Runtime string
	DisableNetwork bool
}

// ContainerServiceInterface is the interface for container servicers
type ContainerServiceInterface interface{
	Run(runContainerOptions RunContainerOptions) (stdOut io.Reader, stdErr io.Reader, err error)
}

