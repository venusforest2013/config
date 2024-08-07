package application

import (
	"testing"
)

var (
	tapp    *Application
	tc      *Controller
	trpc    = &testRPC{}
	tengine = &testEngine{}
)

func TestNewController(t *testing.T) {
	filename := "conf/controller_test.conf"
	cfg, err := newConfigFromTomlFile(filename)
	if err != nil {
		t.Fatal(err)
	}
	tapp = newApp()
	tc = tapp.Controller

	if err := tc.Register(trpc); err != nil {
		t.Fatal(err)
	}
	if err := tc.Register(tengine); err != nil {
		t.Fatal(err)
	}
	if err := tapp.configure(cfg); err != nil {
		t.Fatal(err)
	}
}

func TestController_Module(t *testing.T) {
	if m := tc.Module(trpc.Name()); m != trpc {
		t.Fatal(m)
	}
	if m := tc.Modules(); len(m) != 2 {
		t.Fatal(m)
	}
	if err := tc.start(); err != nil {
		t.Fatal(err)
	}

	if m := tapp.EnabledModules(); len(m) != 1 || m[0] != "test-rpc" {
		t.Fatal(m)
	}
}
