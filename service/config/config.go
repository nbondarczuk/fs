package config

import (
	"fmt"
	"flag"
	"io/ioutil"
)

var (
	Setup             *SetupValueSet
	configFileNamePtr *string
	listenerOnlyMode  bool
)

// checkCmdLineArgs detects usage of cmd line args like -v (version info)
// or -h (default flags). They cause immediate exit but there is a printout
// on the screen with some key info.
func checkCmdLineArgs() {
	v := flag.Bool("v", false, "version info")
	listener := flag.Bool("listener", false, "listener only mode")
	configFileNamePtr = flag.String("config", DEFAULT_CONFIG_FILE_NAME, "config file")
	flag.Parse()

	if *v {
		printVersionInfoAndExit()
	}

	if *listener {
		listenerOnlyMode = true
	}
}

// Init gets the contents of file and uses it to make a config
// It may panic. Handling of it is not a required way.
func Init(gitCommitHash, builtAt, builtBy, builtOn string) {
	setVersionInfo(gitCommitHash, builtAt, builtBy, builtOn)
	checkCmdLineArgs()

	input, err := LoadConfigYamlFromFile(*configFileNamePtr)
	if err != nil {
		panic(err)
	}

	Setup, err = NewSetupValueSet(input)
	if err != nil {
		panic(err)
	}

	Setup.ConfigFileName = *configFileNamePtr
	Setup.ListenerOnlyMode = listenerOnlyMode
	InitLogger(Setup.LogFormat, Setup.LogLevel)
	Setup.LogSetup()
}

// LoadConfigYamlFromFile gets the contents of the config file
func LoadConfigYamlFromFile(fileName string) ([]byte, error) {
	input, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Error opening config file: %s %s", fileName, err)
	}

	return input, nil
}
