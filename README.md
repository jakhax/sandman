# Code Sandbox

## About
Execute and test code of various languages within a sandbox runtime that provides a virtualized container environment.
When code is run it is executed within a docker environment using a [gVisor](https://gvisor.dev/) as the container runtime in order to execute unsafe code in a sandbox.

All execution is done within docker by a compiled go executable within each container that manages the code execution for that specific language environment and returns the result via stdout/ function call result.

**This project was inspired by the now obsolete [codewars cli runner](https://github.com/Codewars/codewars-runner-cli) written in javascript, most of the docker execution environments, test frameworks and docs use components/code from codewars runner.**

## Security 
 **Containers are not contained**
 - Containers by default are not a secure environment to execute untrusted code, an application running in your container can still exploit a vulnerability in the kernel.
 - This why we use [gVisor](https://github.com/google/gvisor) which provides a sanbox runtime.
 - Other ways [Kata Containers](https://katacontainers.io/) and [Firecracker](https://aws.amazon.com/blogs/aws/firecracker-lightweight-virtualization-for-serverless-computing/)
 - For more info about security:
    - [Sandboxing your Containers with gVisor (video)](https://www.youtube.com/watch?v=kxUZ4lVFuVo)
    - [Is docker as a sandbox secure](https://security.stackexchange.com/questions/107850/docker-as-a-sandbox-for-untrusted-code)
    - [Open Sourcing gVisor](https://cloud.google.com/blog/products/gcp/open-sourcing-gvisor-a-sandboxed-container-runtime)
    - [Enemy within; Running untrusted code with gVisor (video)](https://www.youtube.com/watch?v=1Ib-rfSzDuM)
- This project is still a work in progress and is by no means production ready.

## Setup
### Requirements
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Golang]()
### Build executable
```shell
go build -o sandman
```

### Build docker images
- To build all images
```shell
docker-compose build
```
- To build a specific image,find the name of the image in `docker-compose.yaml` for example `python` image
```shell
docker-compose build python-runner
```
### Add gVisor 
- [gVisor](https://gvisor.dev/) is an open-source, OCI-compatible sandbox runtime that provides a virtualized container environment. It runs containers with a new user-space kernel, delivering a low overhead container security solution for high-density applications.
- To install gVisor and use its container runtime see installation instructions [here](https://gvisor.dev/docs/user_guide/quick_start/docker/)


## Basic Usage

### General usage 
```shell
Usage:
  Sandman [flags]
  Sandman [command]

Available Commands:
  help        Help about any command
  run         run code
  run_json    run from json input inside container
  test_       test 
  test_server 

Flags:
  -h, --help   help for Sandman

Use "Sandman [command] --help" for more information about a command.
```
### Running code in sandbox via CLI
- To run code  use the `run` flag `./sandman run`

```shell
Usage:
  Sandman run [flags]

Flags:
  -c, --code string             code to run
  -C, --cpu int                 cpu limit
  -f, --fixture string          Test fixture code to test with
  -h, --help                    help for run
  -l, --language string         The language to execute the code in
  -M, --memory int              memory limit
  -F, --output_format string    Output format, options are 'default' and 'json'
      --sandbox                 environment to run code in, if set will execute code in sandbox
  -s, --setup_code string       Setup code to be used for executing the code
  -S, --shell string            An optional shell script which will be ran within the sandbox environment before the code is executed
      --stdout                  Output stdout and stderr
  -t, --test_framework string   Test framework to use
  -T, --timeout int             The timeout to be used for running the code. If not specified a language specific default will be used
```

For example to run a simple python script which would output `2`:

```shell
./sandman run -c 'print(1+1)' -l 'python' --sandbox -T 200000 --stdout
```

Because everything runs inside of Docker, you would normally not run directly from your host but instead via a Docker run command. To do this use the `--sandbox` flag.

Or you could bash directly into the container
```shell
# direct Docker call:
docker run --rm -it --entrypoint bash sandman/python-runner

# alternatively you can use the provided Docker Compose configuration:
docker-compose run python-runner
```

Or you could choose to execute the code outside of Docker by creating a container that will remove itself after it executes:

```shell
# direct Docker call:
docker run --rm sandman/python-runner run -c 'print(1+1)' -l 'python' -T 200000
docker run --rm codewars/ruby-runner run -l ruby -c "puts 'I ran inside of Docker using Ruby'"

# alternatively you can use the provided Docker Compose configuration:
docker-compose run python-runner -c 'print(1+1)' -l 'python' -T 200000
```

### Run code via script
```golang
```

## Language Support Status

The following languages are currently supported, i'll be adding support for more languages and their respective test frameworks.

| Language       | Version       | Basic Run  | Project Mode | Test Integration | Docker Image  | Notes |
|----------------|---------------|------------|--------------|------------------|---------------|-------------------------------------------------------------------------|
| Go             | 1.10.4           | ✓          |   WIP           | ginkgo           | go-runner |              |
| Python         | 3.6           | ✓          |  WIP            | cw-2, unittest   | python-runner |              |                                                                    
