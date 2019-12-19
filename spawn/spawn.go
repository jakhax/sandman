package spawn
import (
	"io"
	"io/ioutil"
	"bytes"
	"os/exec"
	"time"
	"fmt"
)

// StdIO custom stdout/err with max buffer len
type StdIO struct {
	MaxBuffLen int 
	Buff bytes.Buffer
}

func (b *StdIO) Write(p []byte) (n int, err error){
	
	if b.MaxBuffLen !=0 && b.Buff.Len() + len(p) > b.MaxBuffLen {
		fmt.Println(string(p))
		err =  MaxBufferError{Message: "Max Buffer reached: Too much information has been written to stdio"}
		// sometimes this results to a BrokenPipeError
		return
	}
	
	return b.Buff.Write(p)
}

// NewStdIO returns StdIO
func NewStdIO(maxBufferSize int) (stdio *StdIO){
	var buf bytes.Buffer
	return &StdIO{
		MaxBuffLen: maxBufferSize,
		Buff: buf,
	}
}

// Opt opt
type Opt struct {
	MaxBufferSize int
	Timeout int 
	Dir string
	Env []string

}
// TimeoutError error
type TimeoutError struct{

}

func (e TimeoutError) Error() string{
	return "Process timed out";
}

// MaxBufferError error
type MaxBufferError struct{
	Message string
}

func (e MaxBufferError) Error() string{
	return e.Message;
}

// Spwan spwans a timed process
func Spwan(opt *Opt, name string, args []string, stdin io.Reader)(stdout,stderr io.Reader, err error){
	
	cmd := exec.Command(name,args...);
	

	stdoutBuff :=  NewStdIO(opt.MaxBufferSize)
	stderrBuff := NewStdIO(opt.MaxBufferSize)

	cmd.Stderr = stderrBuff
	cmd.Stdout = stdoutBuff
	
	err = cmd.Start();
	if err != nil{
		return
	}
	var timeoutDuration time.Duration;
	if opt.Timeout != 0{
		timeoutDuration = time.Millisecond*time.Duration(opt.Timeout);
	}else{
		timeoutDuration = time.Second*time.Duration(60);
	}

	timeout := time.After(timeoutDuration);

	done := make(chan error);
	go func(){
		done <- cmd.Wait();
	}()

	select{
		case cmdErr := <-done:
			stdout = &stdoutBuff.Buff
			stderr = &stderrBuff.Buff
			if cmdErr !=  nil{
				if exitErr,ok := cmdErr.(*exec.ExitError); ok{
					err = exitErr
					return
				}
				if maxBufferError,ok := cmdErr.(MaxBufferError); ok{
					err = maxBufferError
					return
				}
				return
			}

		case <- timeout:
			err = cmd.Process.Kill();
			if(err!=nil){
				return;
			}
			err = TimeoutError{};
	}
	return
}

// RunShell run a shell command
func RunShell(opt *Opt, shell []byte) (stdout,stderr io.Reader, err error){
	shellScript :=  string(shell)
	if err !=  nil{
		return
	}
	return Spwan(opt,"echo",[]string{shellScript,"|","bash","-"},nil)
}

//RunShellFromFile runs a shell script
func RunShellFromFile(opt *Opt, file string) (stdout,stderr io.Reader, err error){
	shell,err := ioutil.ReadFile(file)
	if err !=  nil{
		return
	}
	return RunShell(opt, shell);
}