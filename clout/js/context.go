package js

import (
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/dop251/goja"
	"github.com/op/go-logging"
	"github.com/tebeka/atexit"
	"github.com/tliron/kutil/terminal"
	urlpkg "github.com/tliron/kutil/url"
	cloutpkg "github.com/tliron/puccini/clout"
)

//
// Context
//

type Context struct {
	Arguments       map[string]string
	Quiet           bool
	Format          string
	Strict          bool
	AllowTimestamps bool
	Pretty          bool
	Output          string
	Log             *Log
	Stdout          io.Writer
	Stderr          io.Writer
	Stdin           io.Writer
	URLContext      *urlpkg.Context

	programCache sync.Map
}

func NewContext(name string, logger *logging.Logger, arguments map[string]string, quiet bool, format string, strict bool, allowTimestamps bool, pretty bool, output string, urlContext *urlpkg.Context) *Context {
	if arguments == nil {
		arguments = make(map[string]string)
	}

	return &Context{
		Arguments:       arguments,
		Quiet:           quiet,
		Format:          format,
		Strict:          strict,
		AllowTimestamps: allowTimestamps,
		Pretty:          pretty,
		Output:          output,
		Log:             NewLog(logger, name),
		Stdout:          terminal.Stdout,
		Stderr:          terminal.Stderr,
		Stdin:           os.Stdin,
		URLContext:      urlContext,
	}
}

func (self *Context) NewCloutRuntime(clout *cloutpkg.Clout, apis map[string]interface{}) *goja.Runtime {
	runtime := goja.New()
	runtime.SetFieldNameMapper(mapper)

	runtime.Set("puccini", self.NewPucciniAPI())

	runtime.Set("clout", self.NewCloutAPI(clout, runtime))

	for name, api := range apis {
		runtime.Set(name, api)
	}

	return runtime
}

func (self *Context) GetProgram(name string, scriptlet string) (*goja.Program, error) {
	p, ok := self.programCache.Load(scriptlet)
	if !ok {
		program, err := goja.Compile(name, scriptlet, true)
		if err != nil {
			return nil, err
		}
		p, _ = self.programCache.LoadOrStore(scriptlet, program)
	}

	return p.(*goja.Program), nil
}

func (self *Context) Exec(clout *cloutpkg.Clout, scriptletName string, apis map[string]interface{}) error {
	scriptlet, err := GetScriptlet(scriptletName, clout)
	if err != nil {
		return err
	}

	program, err := self.GetProgram(scriptletName, scriptlet)
	if err != nil {
		return err
	}

	runtime := self.NewCloutRuntime(clout, apis)

	_, err = runtime.RunProgram(program)
	return UnwrapException(err)
}

func (self *Context) Fail(message string) {
	if !self.Quiet {
		fmt.Fprintln(self.Stderr, terminal.ColorError(message))
	}
	atexit.Exit(1)
}

func (self *Context) Failf(format string, args ...interface{}) {
	self.Fail(fmt.Sprintf(format, args...))
}

func (self *Context) FailOnError(err error) {
	if err != nil {
		self.Fail(err.Error())
	}
}
