package main

import (
	"io"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

type StratumMinerConfig struct {
	Url      string
	Username string
	Password string
	SumIntv  string
	Threads  int
}

func main() {
	var url, username, password, loglevel, logfile string
	var threads int
	pflag.StringVarP(&url, "url", "o", "erg.uupool.cn:17788", "stratum pool url")
	pflag.StringVarP(&username, "username", "u", "3WxsCRcd9e2fGDHgxk1ptn3stiXPG4p2fi6im4VzQgfQJuS9YSN8.worker1", "username")
	pflag.StringVarP(&password, "password", "x", "x", "password")
	pflag.StringVarP(&loglevel, "loglevel", "l", "debug", "log level: info, debug, trace")
	pflag.StringVarP(&logfile, "logfile", "f", "debug.log", "logfile path")
	pflag.IntVarP(&threads, "threads", "t", runtime.NumCPU(), "threads")
	pflag.Parse()

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	})
	if l, err := logrus.ParseLevel(loglevel); err != nil {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(l)
	}

	if logfile == "" {
		logrus.Warningf("Ignore logging to file")
	}
	ljack := &lumberjack.Logger{
		Filename: logfile,
	}
	mWriter := io.MultiWriter(os.Stdout, ljack)
	logrus.SetOutput(mWriter)

	cfg := &StratumMinerConfig{
		Url:      url,
		Username: username,
		Password: password,
		SumIntv:  "10s",
		Threads:  threads,
	}
	m := NewMiner(cfg)
	m.Mine()
}
