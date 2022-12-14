package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SetupValueSet is a ready to use, parsed structure, contrary to raw config from
// yaml or env variables
type SetupValueSet struct {
	ConfigFileName        string
	ListenerOnlyMode      bool
	LogLevel              string
	LogFormat             string
	UseFileStore          string
	UseListener           bool
	ObjectNamespaceModule uint
	CheckSumType          string
	ServerIPAddress       string
	ServerPort            string
	SQLMaxIdleConns       int
	SQLMaxOpenConns       int
	SQLMaxLifetime        time.Duration
	PresignDurMins        time.Duration
}

// LogSetup show initial start info with the setup of env
func (s *SetupValueSet) LogSetup() {
	Log.Info("Starting")

	logVersionInfo()

	// Config info
	Log.Info("Config setup")
	Log.Info("        Config file name: " + s.ConfigFileName)
	Log.Info("       Listener Only Mode: " + strconv.FormatBool(s.ListenerOnlyMode))

	// Logging info
	Log.Info("               Log Level: " + s.LogLevel)
	Log.Info("              Log Format: " + s.LogFormat)

	// File store paramters
	Log.Info("          Use File Store: " + s.UseFileStore)
	Log.Info("            Use Listener: " + strconv.FormatBool(s.UseListener))
	Log.Info("           Checksum Type: " + s.CheckSumType)
	
	// Server parameters
	Log.Info("    HTTP ServerIPAddress: " + s.ServerIPAddress)
	Log.Info("         HTTP ServerPort: " + s.ServerPort)

	// SQL connection options
	Log.Info("         SQLMaxIdleConns: " + fmt.Sprintf("%d", s.SQLMaxIdleConns))
	Log.Info("         SQLMaxOpenConns: " + fmt.Sprintf("%d", s.SQLMaxOpenConns))
	Log.Info("          SQLMaxLifetime: " + fmt.Sprintf("%d hours", s.SQLMaxLifetime))

	// Postgres credentials
	Log.Info("           POSTGRES_USER: " + os.Getenv("POSTGRES_USER"))
	Log.Info("           POSTGRES_PASS: " + s.Anonym(os.Getenv("POSTGRES_PASS")))
	Log.Info("         POSTGRES_DBNAME: " + os.Getenv("POSTGRES_DBNAME"))
	Log.Info("           POSTGRES_HOST: " + os.Getenv("POSTGRES_HOST"))
	Log.Info("           POSTGRES_PORT: " + os.Getenv("POSTGRES_PORT"))
}

// NewSetup create a new configuration
func NewSetupValueSet(input []byte) (*SetupValueSet, error) {
	s := &SetupValueSet{}
	err := s.load(input)
	if err != nil {
		return nil, fmt.Errorf("Error loading config: %s", err)
	}

	return s, err
}

// load gets config from file, env vars (if found) and command line (if used)
func (s *SetupValueSet) load(input []byte) error {
	s.initDefaultValues()
	err := s.loadFromYaml(input)
	if err != nil {
		return fmt.Errorf("Error loading yaml file: %s", err)
	}

	err = s.setEnvValues()
	if err != nil {
		return fmt.Errorf("Invalid config parameter: %s", err)
	}

	return nil
}

// setInitConfig initializes the config with initial profile which must be not empty
func (s *SetupValueSet) initDefaultValues() {
	s.LogLevel = DEFAULT_LOG_LEVEL
	s.LogFormat = DEFAULT_LOG_FORMAT
	s.UseFileStore = DEFAULT_FILE_STORE
	s.CheckSumType = DEFAULT_CHECKSUM_TYPE
	s.ObjectNamespaceModule = DEFAULT_OBJECT_NAMESPACE_MODULE
	s.PresignDurMins = DEFAULT_PRESIGNED_DUR_MINS
	s.ServerIPAddress = DEFAULT_IP_ADDRESS
	s.ServerPort = DEFAULT_PORT
	s.SQLMaxIdleConns = DEFAULT_SQL_MAX_IDLE_CONNS
	s.SQLMaxOpenConns = DEFAULT_SQL_MAX_OPEN_CONNS
	s.SQLMaxLifetime = time.Hour * DEFAULT_SQL_MAX_LIFETIME
}

// setEnvValues uses env vars to set specific confi flags
func (s *SetupValueSet) setEnvValues() error {
	var val string

	val = os.Getenv("LOG_LEVEL")
	if val != "" {
		switch val {
		case "debug":
		case "info":
			s.LogLevel = val
		default:
			return fmt.Errorf("Invalid log level: %s, expected: info, debug", val)
		}
	}

	val = os.Getenv("LOG_FORMAT")
	if val != "" {
		switch val {
		case "json":
		case "text":
			s.LogFormat = val
		default:
			return fmt.Errorf("Invalid log level: %s, expected: kjson, text", val)
		}
	}

	val = os.Getenv("USE_LISTENER")
	if val != "" {
		if val == "true" {
			s.UseListener = true
		}
	}

	val = os.Getenv("USE_LISTENER_ONLY_MODE")
	if val != "" {
		if val == "true" {
			s.ListenerOnlyMode = true
		}
	}

	val = os.Getenv("USE_OBJECT_NAMESPACE_MODULE")
	if val != "" {
		valuint64, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "USE_OBJECT_NAMESPACE_MODULE", val)
		}

		s.ObjectNamespaceModule = uint(valuint64)
	}

	val = os.Getenv("HTTP_ADDRESS")
	if val != "" {
		s.ServerIPAddress = val
	}

	val = os.Getenv("HTTP_PORT")
	if val != "" {
		s.ServerPort = val
	}

	var err error
	
	val = os.Getenv("SQL_MAX_IDLE_CONNS")
	if val != "" {
		s.SQLMaxIdleConns, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_IDLE_CONNS", val)
		}
	}

	val = os.Getenv("SQL_MAX_OPEN_CONNS")
	if val != "" {
		s.SQLMaxOpenConns, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_OPEN_CONNS", val)
		}
	}

	val = os.Getenv("SQL_MAX_LIFETIME")
	if val != "" {
		var valint int
		valint, err = strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("Invalid env variable %s value: %s", "SQL_MAX_LIFETIME", val)
		}
		s.SQLMaxLifetime = time.Hour * time.Duration(valint)
	}

	s.logFileStoreInfo()

	return nil
}

func (s *SetupValueSet) logFileStoreInfo() {
	val := os.Getenv("USE_FILE_STORE")
	if val != "" {
		s.UseFileStore = val
	}
}

// loadFromYamlFile loads the config.yaml file overriding default config
func (s *SetupValueSet) loadFromYaml(input []byte) error {
	var doc Document
	err := yaml.Unmarshal(input, &doc)
	if err != nil {
		return err
	}

	for _, logger := range doc.Loggers {
		s.setEnvVars(logger.Kind, logger.Env)
	}

	for _, server := range doc.Servers {
		s.setEnvVars(server.Kind, server.Env)
	}

	for _, sqloption := range doc.SQLOptions {
		s.setEnvVars(sqloption.Kind, sqloption.Env)
	}

	for _, backend := range doc.Backends {
		s.setEnvVars(backend.Kind, backend.Env)
	}

	for _, store := range doc.Stores {
		s.setEnvVars(store.Kind, store.Env)
	}

	return nil
}

// setEnvVars overrides the default values from config with the env
func (s *SetupValueSet) setEnvVars(kind string, env map[string]string) {
	for key, envval := range env {
		envvar := fmt.Sprintf("%s_%s", strings.ToUpper(kind), strings.ToUpper(key))
		if os.Getenv(envvar) == "" {
			flgval := checkCmdLineUsage(strings.ToLower(key))
			if flgval == "" {
				os.Setenv(envvar, envval)
			} else {
				os.Setenv(envvar, flgval)
			}
		}
	}
}

// checkFlagsUsage TBD. it looks for possible use of an option in the command line
func checkCmdLineUsage(flg string) (val string) {
	//flag.StringVar(&val, flg, "", "")
	//flag.Parse()
	return
}

// anonymize reveals secrets in debug or trace mode
func (s *SetupValueSet) Anonym(str string) string {
	if s.LogLevel == "debug" {
		return str
	}

	return "[...]"
}
