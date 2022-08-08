package config

import (
	"bufio"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/zhaolion/gostack/util/xerr"
)

type AppConfig struct {
	AppName  string         `json:"app_name" yaml:"app_name" toml:"app_name"`
	AppEnv   string         `json:"app_env" yaml:"app_env" toml:"app_env" env:"APP_ENV"`
	Debug    bool           `json:"debug" yaml:"debug" toml:"debug" env:"APP_DEBUG"`
	Database DatabaseConfig `json:"database" yaml:"database" toml:"database"`
}

type DatabaseConfig struct {
	DSN string `json:"dsn" yaml:"dsn" toml:"dsn" env:"DB_DSN"`
}

type CSVConfig struct {
	AppName          string            `json:"app_name" yaml:"app_name" toml:"app_name"`
	AppEnv           string            `json:"app_env" yaml:"app_env" toml:"app_env" env:"APP_ENV"`
	Debug            bool              `json:"debug" yaml:"debug" toml:"debug" env:"APP_DEBUG"`
	Database         DatabaseConfig    `json:"database" yaml:"database" toml:"database"`
	LessonExtensions []LessonExtension `json:"-" yaml:"-" toml:"-" file:"test-csv.csv"`
}

type LessonExtension struct {
	ID string `csv:"id"`
	// 封面图片, optional
	CoverURL string `csv:"cover_url"`
	// 预估时长, optional
	DurationSeconds int32 `csv:"duration_seconds"`
	// 课程等级难度, optional
	Grade string `csv:"grade"`
}

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type Suite struct {
	suite.Suite

	backupEnvs []string

	JSONFile string
	EnvFile  string
	YAMLFile string
	TOMLFile string
	CSVFile  string
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *Suite) SetupSuite() {
	suite.JSONFile = "test-json.json"
	initConfig(suite.JSONFile, jsonExample)

	suite.EnvFile = "test-env.json"
	initConfig(suite.EnvFile, jsonExample)

	suite.YAMLFile = "test-yaml.yml"
	initConfig(suite.YAMLFile, yamlExample)

	suite.TOMLFile = "test-toml.toml"
	initConfig(suite.TOMLFile, tomlExample)

	suite.CSVFile = "test-csv.csv"
	initConfig(suite.CSVFile, csvExample)
}

func (suite *Suite) TearDownSuite() {
	_ = os.Remove(suite.JSONFile)
	_ = os.Remove(suite.EnvFile)
	_ = os.Remove(suite.YAMLFile)
	_ = os.Remove(suite.TOMLFile)
	_ = os.Remove(suite.CSVFile)
}

// The SetupTest method will be run before every test in the suite.
func (suite *Suite) SetupTest() {
	// protected env keys
	suite.backupEnvs = os.Environ()
	os.Clearenv()
}

// The TearDownTest method will be run after every test in the suite.
func (suite *Suite) TearDownTest() {
	for _, pair := range suite.backupEnvs {
		kvs := strings.Split(pair, "=")

		if len(kvs) > 2 {
			os.Setenv(kvs[0], kvs[1])
		}
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func initConfig(file, config string) {
	Reset()
	outputFile, outputError := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0666)
	if outputError != nil {
		panic(xerr.New("an error occurred with file opening or creation"))
	}
	defer func() {
		_ = outputFile.Close()
	}()

	outputWriter := bufio.NewWriter(outputFile)
	_, _ = outputWriter.WriteString(config)
	_ = outputWriter.Flush()
}

var jsonExample = `
{
	"app_name": "test-app",
	"app_env": "test",
	"debug": true,
	"database": {
	  "dsn": "root@tcp(localhost:3306)/test"
	}
  }
`

var yamlExample = `
app_name: test-app
app_env: test
debug: true
database:
  dsn: "root@tcp(localhost:3306)/test"
`

var tomlExample = `
app_name = "test-app"
app_env = "test"
debug = true

[database]
dsn = "root@tcp(localhost:3306)/test"
`

var csvExample = `id,cover_url,duration_seconds,grade
100000006c111bea29e0b02c,https://z-b.llscdn.com/zhenghe-yy-2c/8dcbf5589e6b3b62c9b86a04494a288d8bed9899.jpg,180,ELEMENTARY
1000000077174595134a18e0,https://z-b.llscdn.com/zhenghe-yy-2c/98c791d8f36e3c650883b9bcafe9b1e7481d1ce9.jpeg,380,INTERMEDIATE
100000002a630c13856f4b95,https://z-b.llscdn.com/zhenghe-yy-2c/4585be7f9160bc6bda707cf10899a9c3879ef979.jpg,480,ADVANCED`
