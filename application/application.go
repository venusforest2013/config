package application

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync"
	"syscall"

	"github.com/venusforest2013/config"
)

const (
	Key = "__APPLICATION"
)

var (
	Version   = "0.0.0"
	Revision  = "master"
	BuildTime = "None"

	defaultMaxThreads     = int64(10000)
	defaultMaxProcesses   = int64(0)
	defaultMaxStackMBytes = int64(128)
	defaultGCPercent      = int64(100)

	appFlagsDaemon = "__APPLICATION_DAEMON"
	appFileMode    = os.FileMode(0644)

	app = newApp()
)

func New() *Application {
	return app
}

func Context() *Application {
	return app
}

// 注册模块，时机：模块包 init()
func Register(module Module) {
	if err := app.Register(module); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}

func Configure(cfg *Config) *Application {
	if err := app.configure(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	return app
}

func GetModule(name string) Module {
	if m := app.Controller.Module(name); m != nil {
		return m
	}
	log.Panicf("Module '%s' not found", name)
	return nil
}

// Config holds configuration for application settings.
type Config = config.Value

type Application struct {
	*Controller
	ctx          context.Context
	name         string
	config       *Config
	logger       *Logger
	signals      chan os.Signal
	once         *sync.Once
	pidFile      string
	configSource string
}

func newApp() *Application {
	app := &Application{
		name:    cmd,
		logger:  NewLogger(os.Stderr, true),
		signals: make(chan os.Signal, 1),
		once:    &sync.Once{},
	}
	app.ctx = context.WithValue(context.Background(), Key, app)
	app.Controller = NewController(app.ctx)
	return app
}

func (app *Application) Name() string {
	return app.name
}

func (app *Application) Config() *Config {
	return app.config
}

func (app *Application) Logger() *Logger {
	return app.logger
}

func (app *Application) EnabledModules() []string {
	modules := []string{}
	for _, m := range app.Controller.Modules() {
		if app.Controller.Enabled(m.Name()) {
			modules = append(modules, m.Name())
		}
	}
	return modules
}

func (app *Application) Start() error {
	var err error
	app.once.Do(func() {
		err = app.start()
	})
	return err
}

func (app *Application) configure(cfg *Config) error {
	var err error
	defer func() {
		if err != nil {
			app.logger.Printf("failed to configure %s: %v", app.name, err)
		}
	}()

	// Runtime
	debug.SetMaxThreads(int(cfg.Int64("max_threads", defaultMaxThreads)))
	runtime.GOMAXPROCS(int(cfg.Int64("max_processes", defaultMaxProcesses)))
	debug.SetMaxStack((1024 * 1024) * int(cfg.Int64("max_stack_mb", defaultMaxStackMBytes)))
	debug.SetGCPercent(int(cfg.Int64("gc_percent", defaultGCPercent)))

	// Logger
	logFile := cfg.Str("log_file", "")
	if logFile != "" {
		logger, err := NewFileLogger(logFile, true)
		if err != nil {
			return err
		}
		app.logger = logger
	}

	if err = app.Controller.configure(cfg); err != nil {
		return err
	}

	app.pidFile = cfg.Str("pid_file", fmt.Sprintf("./%s.pid", cmd))
	app.config = cfg
	return nil
}

func (app *Application) start() error {

	if app.config == nil {
		if err := configure(app); err != nil {
			return err
		}
	}

	if err := app.daemonize(); err != nil {
		return err
	}

	app.logger.Printf("starting %s ......", app.name)

	// pid 文件加锁，防止重复启动
	if err := app.lock(); err != nil {
		return fmt.Errorf("failed to start server: %v", err)
	}

	if app.configSource != "" {
		app.logger.Printf("loaded configuration source => %s", app.configSource)
	}

	// 启动模块
	if err := app.Controller.start(); err != nil {
		return fmt.Errorf("failed to start modules: %v", err)
	}
	app.logger.Printf("loaded modules => %v", app.EnabledModules())

	app.logger.Printf("%s started successfully", app.name)
	return app.wait()
}

func (app *Application) shutdown() error {
	var err error
	defer func() {
		if err != nil {
			app.logger.Printf("failed to shutdown %s: %v", app.name, err)
		} else {
			app.logger.Printf("%s shutdown successfully", app.name)
		}
	}()

	modules := app.EnabledModules()
	if err = app.Controller.shutdown(); err != nil {
		return err
	}
	app.logger.Printf("shutdown modules => %v", modules)

	if err = os.Remove(app.pidFile); err != nil {
		return err
	}
	return nil
}

func (app *Application) kill(signal string) error {
	if signal == "" {
		return nil
	}
	cpidstr, err := ioutil.ReadFile(app.pidFile)
	if err != nil {
		return err
	}
	cpid, err := strconv.Atoi(string(cpidstr))
	if err != nil {
		return err
	}

	if err = kill(cpid, signal); err == nil {
		os.Exit(0)
	}
	return err
}

func (app *Application) daemonize() error {
	daemon := app.config.Bool("daemon", false)
	if !daemon || os.Getenv(appFlagsDaemon) == "1" {
		return nil
	}
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return err
	}
	dir, _ := os.Getwd()
	env := append(os.Environ(), appFlagsDaemon+"=1")
	cmd := &exec.Cmd{
		Path:         path,
		Args:         os.Args,
		Env:          env,
		Dir:          dir,
		Stdin:        nil,
		Stdout:       nil,
		Stderr:       nil,
		ExtraFiles:   nil,
		SysProcAttr:  &syscall.SysProcAttr{Setsid: true},
		Process:      nil,
		ProcessState: nil,
	}
	err = cmd.Start()
	if err == nil {
		os.Exit(0)
	}
	return err
}

func (app *Application) lock() error {
	filename := app.pidFile
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, appFileMode)
	if err != nil {
		return err
	}
	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		return err
	}
	file.Truncate(0)
	_, err = file.Write([]byte(strconv.Itoa(os.Getpid())))
	return err
}

func (app *Application) wait() error {
	ss := []os.Signal{
		syscall.SIGQUIT, syscall.SIGTERM,
	}
	signal.Notify(app.signals, ss...)
	defer signal.Stop(app.signals)
	for {
		select {
		case sig := <-app.signals:
			switch sig {
			case syscall.SIGQUIT, syscall.SIGTERM:
				app.logger.Printf("caught %s, shutting down ......\n", sig)
				return app.shutdown()
			default:
				app.logger.Printf("unsupported signal %v", sig)
			}
		}
	}
	return nil
}
