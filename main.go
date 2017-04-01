package main

import (
	"github.com/cdelorme/glog"
	"github.com/cdelorme/gonf"
)

const defaultPackageName = "com.example"

func main() {
	logger := &glog.Logger{}
	app := &appinator{Package: defaultPackageName}

	conf := gonf.Config{}
	conf.Target(app)
	conf.Description("simple osx application bundler")
	conf.Add("name", "override for app name", "", "-n", "--name")
	conf.Add("app", "executable path", "", "-a", "--app")
	conf.Add("package", "package name (default `com.example`)", "", "-p", "--package")
	conf.Add("icon", "icon path", "", "-i", "--icon")
	conf.Add("resources", "resources path", "", "-r", "--resources")
	conf.Add("frameworks", "frameworks path", "", "-f", "--frameworks")
	conf.Add("docs", "documentation path", "", "--docs")
	conf.Add("debug", "enable debug, use symlink", "", "-d", "--debug")
	conf.Load()

	logger.Debug("%#v", app)

	if err := app.Build(); err != nil {
		logger.Error(err.Error())
	}
}
