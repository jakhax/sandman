package utils;


// ErrRunningContainer when a container fails to run
type ErrRunningContainer struct{
	Message string
}

func (e ErrRunningContainer) Error() string{
	return e.Message;
}

// ErrContainerTimeout when a container timesout
type ErrContainerTimeout struct{
	Message string
}

func (e ErrContainerTimeout) Error() string{
	return e.Message;
}

// InvalidSandBoxOpt for sandbox options
type InvalidSandBoxOpt struct{
	Message string
}

func (e InvalidSandBoxOpt) Error() string{
	return e.Message;
}