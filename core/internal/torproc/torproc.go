// Package torproc runs the downloaded tor binary for the tor outbound, always as the invoking user.
package torproc

import (
	"context"
	"errors"
	"net"
	"path/filepath"
	"runtime"

	rawlineprocess "RawlineCore/internal/process"

	"github.com/cretz/bine/process"
)

// Relative to the core's working directory, the config folder.
const (
	Folder        = "tor"
	DataDirectory = Folder + "/data"
)

func ExecutableName() string {
	if runtime.GOOS == "windows" {
		return "tor.exe"
	}
	return "tor"
}

// Never "": an empty executable_path makes the outbound look tor up on PATH.
func ExecutablePath() string {
	return absolute(filepath.Join(Folder, ExecutableName()))
}

// Absolute, so tor resolves it the same whatever directory it was started in.
func DataDirectoryPath() string {
	return absolute(filepath.FromSlash(DataDirectory))
}

func absolute(path string) string {
	if abs, err := filepath.Abs(path); err == nil {
		return abs
	}
	return path
}

type Creator struct{}

func (Creator) NewProcess(ctx context.Context, executablePath string, args ...string) (process.Process, error) {
	if filepath.Clean(executablePath) != ExecutablePath() {
		return nil, errors.New("Tor must use the executable installed by Rawline")
	}
	return &torProcess{ctx: ctx, path: executablePath, args: args}, nil
}

type torProcess struct {
	ctx    context.Context
	path   string
	args   []string
	exited chan struct{}
	err    error
}

func (p *torProcess) Start() error {
	if p.exited != nil {
		return errors.New("tor already started")
	}
	if err := p.ctx.Err(); err != nil {
		return err
	}
	run, err := rawlineprocess.StartTor(p.path, p.args)
	if err != nil {
		return err
	}
	p.exited = make(chan struct{})
	go func() {
		p.err = run.Wait()
		close(p.exited)
	}()
	go func() {
		select {
		case <-p.ctx.Done():
			_ = run.Kill()
		case <-p.exited:
		}
	}()
	return nil
}

func (p *torProcess) Wait() error {
	if p.exited == nil {
		return errors.New("tor not started")
	}
	<-p.exited
	return p.err
}

func (p *torProcess) EmbeddedControlConn() (net.Conn, error) {
	return nil, process.ErrControlConnUnsupported
}
