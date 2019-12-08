package runners;
import (
	"io"
	"bytes"
	"os/exec"
	"fmt"
	"time"
	"errors"
)

// SpawnOpt opt
type SpawnOpt struct {
	Timeout int 
	Dir string
	Env []string

}
// TimeoutError error
type TimeoutError struct{
	Message string
}

func (e TimeoutError) Error() string{
	return e.Message;
}

// MaxBufferError error
type MaxBufferError struct{
	Message string
}

func (e MaxBufferError) Error() string{
	return e.Message;
}

// Spwan spwans a timed process
func Spwan(opt *SpawnOpt, name string, args []string, stdin io.Reader)(stdout,stderr io.Reader, err error){
	cmd := exec.Command(name,args...);
	stderrPipe,err := cmd.StderrPipe()
	if err != nil{
		
		return
	}
	stdoutPipe,err := cmd.StdoutPipe()
	stdout = stdoutPipe
	stderr = stderrPipe
	if err != nil{
		return
	}

	stdErrChan :=  make(chan StdIO)
	stdOutChan :=  make(chan StdIO)
	stdioErr := make(chan error);
	complete := make(chan bool);

	
	err = cmd.Start();
	if err != nil{
		fmt.Println("errr");
		return
	}
	go ReadStdIOFromPipe(stdoutPipe,stderrPipe,stdOutChan,stdErrChan,complete,stdioErr)

	done := make(chan error);
	go func(){
		done <- cmd.Wait();
		
	}()

	var timeoutDuration time.Duration;
	if opt.Timeout != 0{
		timeoutDuration = time.Second*time.Duration(opt.Timeout);
	}else{
		timeoutDuration = time.Second*time.Duration(15);
	}

	timeout := time.After(timeoutDuration);
	
	select{
		case cmdErr := <-done:
			
			complete <- true;
			
			stdOut := <-stdOutChan;
			stdErr := <-stdErrChan;
			
			stdout = stdOut.Data
			if stdOut.Error != nil{
				err = stdOut.Error
				return
			}
			stderr = stdErr.Data
			if stdErr.Error != nil{
				err = stdErr.Error
				return
			}
			if cmdErr != nil{
				err = cmdErr
				return;
			}
			
		case <- timeout:
			err = cmd.Process.Kill();
			if(err!=nil){
				return;
			}
			complete <- false
			err = errors.New("Process timed out");
		case err = <- stdioErr:
			if err != nil{
				return
			}
		}
	

	return
}

// StdIO to store stdio from child process
type StdIO struct{
	Data io.Reader
	Error error
}

// ReadStdIOFromPipe to read from child process stdout and stderr pipe
func ReadStdIOFromPipe(stdoutPipe,stderrPipe io.ReadCloser, stdout,stderr chan StdIO, ok chan bool, stdioErr chan error){
	reader := func(pipe io.ReadCloser, stdio chan io.Reader, errCh chan error, d string){
		
		KB := 1024;
		MAX_BUFFER := KB * 1500;
		var out []byte;
		buf := make([]byte, 1024, 1024)
		var err error;
		for{
			n,errX := pipe.Read(buf[:]);
			if errX != nil{
				if errX != io.EOF{
					err = errX;
				}
				break;
			}
			if n>0{
				if len(out)+n > MAX_BUFFER{
					err = MaxBufferError{
						Message:"Max Buffer reached: Too much information has been written to stdio",
					};
					break;
				}
				d :=  buf[:n];
				out =  append(out,d...);
			}
		}		
		r := bytes.NewReader(out);
		stdio <- r;
		errCh <- err;
	}
	stdoutChan := make(chan io.Reader)
	errStdoutChan := make(chan error)
	stderrChan :=  make(chan io.Reader)
	errStderrChan := make(chan error)
	go reader(stdoutPipe,stdoutChan,errStdoutChan,"stdout")
	go reader(stderrPipe,stderrChan,errStderrChan,"stderr")
	stdoutR := <- stdoutChan
	stdoutErr := <-errStdoutChan
	if err,ok := stdoutErr.(MaxBufferError); ok{
		stdioErr <- err
		return
	}
	stderrR := <- stderrChan
	stderrErr := <-errStderrChan
	if err,ok := stderrErr.(MaxBufferError); ok{
		stdioErr <- err
		return
	}
	
	complete := <- ok;
	if complete == true{
		stdout <- StdIO{
			Data:stdoutR,
			Error:stdoutErr,
		}
		stderr <- StdIO{
			Data:stderrR,
			Error:stderrErr,
		}
	}
}