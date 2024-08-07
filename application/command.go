package application

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"

	"github.com/venusforest2013/config"
)

var (
	cmd = command()

	helpOption     = false
	testOption     = false
	versionOption  = false
	signalOption   = ""
	filenameOption = ""

	errConfigFileNotFound = errors.New("file not found")
)

func configure(app *Application) error {
	parseFlagSet()
	if helpOption {
		printUsage()
	} else if versionOption {
		printVersion()
	}

	filename, err := filepath.Abs(filenameOption)
	if err != nil {
		return err
	}
	cfg, err := newConfigFromTomlFile(filename)
	if err == nil {
		err = app.configure(cfg)
	}
	if testOption {
		testExit(filename, err)
	} else if err != nil {
		return err
	}
	if signalOption != "" {
		return app.kill(signalOption)
	}
	app.configSource = filename
	return nil
}

func parseFlagSet() *flag.FlagSet {

	set := flag.NewFlagSet(cmd, flag.ExitOnError)
	set.Usage = printUsage
	set.BoolVar(&helpOption, "h", false, "")
	set.BoolVar(&testOption, "t", false, "")
	set.BoolVar(&versionOption, "v", false, "")
	set.StringVar(&signalOption, "s", "", "")
	set.StringVar(&filenameOption, "c", "", "")
	// 过滤 go toolchain
	if !strings.HasPrefix(cmd, "go") &&
		!strings.HasSuffix(cmd, ".test") {
		set.Parse(os.Args[1:])
	}
	return set
}

func command() string {
	path := strings.Split(os.Args[0], "/")
	return path[len(path)-1]
}

func kill(pid int, signal string) (err error) {
	switch signal {
	case "quit":
		err = syscall.Kill(pid, syscall.SIGQUIT)
	default:
		err = errors.New("unsupported command: " + signal)
	}
	return
}

func testExit(filename string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: configuration file '%s' test failed: %v\n", cmd, filename, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s: the configuration file '%s' syntax is ok\n", cmd, filename)
	}
	os.Exit(0)
}

func printUsage() {
	usage := `Usage: ` + cmd + ` [-thv] [-c filename] [-s signal]

options:
    -h                  show this message and exit
    -v                  show version and exit
    -t                  test configuration file and exit
    -c filename         set configuration file
    -s signal           send signal to ` + cmd + `: quit
`
	fmt.Fprintf(os.Stderr, "%s\n", usage)
	os.Exit(0)
}

func printVersion() {
	fmt.Fprintf(os.Stderr, "%s version %s %s/%s\n",
		cmd, Version, runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(os.Stderr, "git revision: %s\n", Revision)
	fmt.Fprintf(os.Stderr, "built in: %s\n", BuildTime)
	os.Exit(0)
}

func newConfigFromTomlFile(filename string) (*Config, error) {
	if filename == "" {
		return nil, errConfigFileNotFound
	}
	toml, err := config.LoadTomlFile(filename)
	if err != nil {
		e := fmt.Errorf("loaded configuration source '%s' => %v", filename, err)
		return nil, e
	}
	return &config.Value{toml}, nil
}
