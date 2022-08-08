package config

import "os"

func (suite *Suite) TestJSON() {
	cfg := &AppConfig{}
	gotErr := Initialize(suite.JSONFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*AppConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

func (suite *Suite) TestENV() {
	_ = os.Setenv("APP_ENV", "development")
	_ = os.Setenv("APP_DEBUG", "false")
	_ = os.Setenv("DB_DSN", "test@tcp(localhost:3306)")

	cfg := &AppConfig{}
	gotErr := Initialize(suite.EnvFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*AppConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("development", got.AppEnv)
	suite.Equal(false, got.Debug)
	suite.Equal("test@tcp(localhost:3306)", got.Database.DSN)
}

func (suite *Suite) TestYAML() {
	cfg := &AppConfig{}
	gotErr := Initialize(suite.YAMLFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*AppConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

func (suite *Suite) TestTOML() {
	cfg := &AppConfig{}
	gotErr := Initialize(suite.TOMLFile, cfg)
	suite.NoError(gotErr)
	got, ok := Get().(*AppConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.Equal("test", got.AppEnv)
	suite.Equal(true, got.Debug)
	suite.Equal("root@tcp(localhost:3306)/test", got.Database.DSN)
}

func (suite *Suite) TestCSV() {
	cfg := &CSVConfig{}
	gotErr := Initialize(suite.JSONFile, cfg)
	suite.Require().NoError(gotErr)
	got, ok := Get().(*CSVConfig)
	suite.True(ok)
	suite.Equal("test-app", got.AppName)
	suite.NotEmpty(got.LessonExtensions)
}
