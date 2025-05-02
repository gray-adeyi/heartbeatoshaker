package internal

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"
	"errors"
		"os/signal"
		"syscall"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile              = "heartbeat.yaml"
	ExitCodeUnableToReadConfigFile = 3
	ExitCodeUnableToWriteConfigFile= 4
	ExitCodeNoTarget               = 5
	Version = "0.1.0"
)

var (
	url        string
	interval   time.Duration
	configFile string
)

func RunCmd(){
	flag.Usage = usage
		if len(os.Args) >= 2 && os.Args[1] == "init"{
			if err :=initConfig(); err != nil{
				os.Stderr.WriteString("unable to create heartbeat config")
				os.Exit(ExitCodeUnableToWriteConfigFile)
			}
			os.Exit(0)
		}
		if len(os.Args) >= 2 && os.Args[1] == "version"{
			fmt.Printf("Heartbeatoshaker Version(%s)\n", Version)
			os.Exit(0)
		}
		config := Config{}
		parseFlags()
		if err :=  DefaultConfigFileReader.Read(&config); err != nil {
			if !errors.Is(err, ErrNoDefaultConfigFile) {
				os.Stderr.WriteString("unable to read your config file\n")
				os.Exit(ExitCodeUnableToReadConfigFile)
			}
			if !urlTargetFromFlagsIsValid() {
				os.Stderr.WriteString("no config file was found and no -url was passed to the cli\n")
				os.Exit(ExitCodeNoTarget)
			}
			setTargetsFromFlags(&config)
		}
		heart := NewHeart(config.Targets)
		heart.Beat()

		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		fmt.Println("Heartbeatoshaker runninng...\npress CTRL+C to cancel")

		for {
			select {
			case <-sigCh:
				go heart.Stop()
			case <-heart.Done:
			os.Exit(0)
			}
		}
}

func parseFlags() {
	flag.StringVar(&url, "url", "", "The target url")
	flag.DurationVar(&interval, "interval", 60*time.Second, "The intervals to ping")
	flag.StringVar(&configFile, "config", DefaultConfigFile, "Alternate path to the hearbeat.yaml file")
	flag.Parse()
}

func urlTargetFromFlagsIsValid() bool {
	return url != "" && interval != 0*time.Second
}

func setTargetsFromFlags(config *Config) {
	config.Targets = []UrlTarget{
		{Url: url, Interval: interval},
	}
}


func initConfig()  error {
	config :=  Config{
		Targets: []UrlTarget{
				{Url: "https://github.com/gray-adeyi/heartbeatoshaker", Interval: 1*time.Second},
				{Url: "https://github.com/gray-adeyi", Interval: 5*time.Minute},
			},
	}
	data, err :=yaml.Marshal(config)
	if err != nil{
		return err
	}
	os.WriteFile(DefaultConfigFile, data,fs.ModeAppend)
	return nil
}


func usage(){
	message := `Heartbeatoshaker Version(%s)

This is a utility that is used to periodically ping a url or multiple urls if a config file is
provided.
	
Usages:
	hbs (Runs heartbeatoshaker with the config file it finds)
	hbs init (This creates a heartbeat.yaml file in the current path that you can configure)
	hbs -url https://github.com/gray-adeyi/heartbeatoshaker  -interval 5m (Pings the url)
	hbs -config my-heartbeat.yml (Runs heartbeatoshaker with the custom config file)
`
	fmt.Print(fmt.Sprintf(message, Version))
}

