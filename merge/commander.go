package merge

import (
	"io"
	"os/exec"
)

type (
	// Commander is a wrapper for exec.Command()
	Commander interface {
		Command(name string, arg ...string) Cmd
	}

	// Cmd is a wrapper for exec.Cmd{}
	Cmd interface {
		StdoutPipe() (io.ReadCloser, error)
		SetEnvironment(env []string)
		StderrPipe() (io.ReadCloser, error)
		Start() error
		Wait() error
	}

	commanderImpl struct {
	}

	cmdImpl struct {
		cmd *exec.Cmd
	}
)

// NewCommander instantiates a new Commander.
func NewCommander() Commander {
	return &commanderImpl{}
}

func (c *commanderImpl) Command(name string, arg ...string) Cmd {
	cmd := exec.Command(name, arg...)
	return &cmdImpl{
		cmd: cmd,
	}
}

func (c *cmdImpl) StdoutPipe() (io.ReadCloser, error) {
	return c.cmd.StdoutPipe()
}

func (c *cmdImpl) SetEnvironment(env []string) {
	c.cmd.Env = env
}

func (c *cmdImpl) StderrPipe() (io.ReadCloser, error) {
	return c.cmd.StderrPipe()
}

func (c *cmdImpl) Start() error {
	return c.cmd.Start()
}

func (c *cmdImpl) Wait() error {
	return c.cmd.Wait()
}
