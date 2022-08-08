package encode

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
)

func TestParamValidate(t *testing.T) {
	unsupported := []interface{}{
		0,
		"",
		[]string{},
		map[string]string{},
		true,
	}
	for _, v := range unsupported {
		if err := Encode(v, "env", getEnv); err == nil {
			t.Error("expected: pointer to struct type required error")
		}
		if err := Encode(&v, "env", getEnv); err == nil {
			t.Error("expected: pointer to struct type required error")
		}
	}
}

func TestEncode(t *testing.T) {
	type AppMode string
	type Basic struct {
		Mode AppMode `env:"struct1_mode"`
		Addr string  `env:"struct1_addr"`
	}
	type struct1 struct {
		Basic
		Name   string    `env:"struct1_name"`  // string
		Hosts  []string  `env:"struct1_hosts"` // slice
		Hosts1 [3]string `env:"struct1_hosts"` // array
		Port   int       `env:"struct1_port"`  // int
		Debug  bool      `env:"struct1_debug"` // bool
		Rate   float64   `env:"struct1_rate"`  // float
		Omit   string    `env:"-"`             // omit
		DB     struct {
			Drive string   `env:"struct1_db_drive"`
			Hosts []string `env:"struct1_db_hosts"`
			Port  uint     `env:"struct1_db_port"`
		}
	}
	envs := map[string]string{
		"struct1_mode":     "test",
		"struct1_addr":     "localhost",
		"struct1_name":     "app",
		"struct1_hosts":    "192.168.33.1,192.168.33.2,192.168.33.3",
		"struct1_port":     "8080",
		"struct1_debug":    "true",
		"struct1_rate":     "0.618",
		"struct1_omit":     "omit",
		"struct1_db_drive": "mysql",
		"struct1_db_hosts": "1,2,3",
		"struct1_db_port":  "3306",
	}
	value := struct1{
		Name:   "app",
		Hosts:  []string{"192.168.33.1", "192.168.33.2", "192.168.33.3"},
		Hosts1: [3]string{"192.168.33.1", "192.168.33.2", "192.168.33.3"},
		Port:   8080,
		Debug:  true,
		Rate:   0.618,
		Omit:   "",
	}
	value.Mode = "test"
	value.Addr = "localhost"
	value.DB.Drive = "mysql"
	value.DB.Hosts = []string{"1", "2", "3"}
	value.DB.Port = 3306

	for k, v := range envs {
		os.Setenv(k, v)
	}
	var s struct1
	err := Encode(&s, "env", getEnv)
	if err != nil {
		t.Errorf("expected: %v, got: %v", nil, err)
	}
	if err != nil {
		return
	}
	eq, err := equal(s, value)
	if err != nil {
		t.Fatal(err)
	}
	if !eq {
		t.Errorf("\nexpected: %+v\ngot: %+v", value, s)
	}
}

func TestEncodeDVA(t *testing.T) {
	type struct2 struct {
		Str       string `dva:"app.env1"`
		Int       int    `dva:"app.env2"`
		Bool      bool   `dva:"app.env3"`
		JSONDummy Dummy  `dva:"dummy.json"`
		YAMLDummy Dummy  `dva:"dummy.yaml"`
	}
	s2 := struct2{}
	dva := mockDva{}
	err := Encode(&s2, "dva", dva.getDva)
	if err != nil {
		t.Fatal(err)
	}
	if s2.Str != "A" {
		t.Errorf("expected: %+v, got: %+v", "A", s2.Str)
	}
	if s2.Int != 1 {
		t.Errorf("expected: %+v, got: %+v", 1, s2.Int)
	}
	if !s2.Bool {
		t.Errorf("expected: %+v, got: %+v", true, s2.Bool)
	}
	if s2.JSONDummy.A != "123" {
		t.Errorf("expected: %+v, got: %+v", "123", s2.JSONDummy.A)
	}
	if s2.YAMLDummy.A != "123" {
		t.Errorf("expected: %+v, got: %+v", "123", s2.YAMLDummy.A)
	}
}

func TestSetBoolVal(t *testing.T) {
	type TC struct {
		Value  string
		HasErr bool
		Bool   bool
	}
	tcs := []TC{
		TC{"", false, false},
		TC{"true", false, true},
		TC{"false", false, false},
		TC{"True", true, false},
		TC{"True", true, false},
		TC{"TRUE", true, false},
		TC{"1", true, false},
	}
	for _, tc := range tcs {
		var b bool
		v := reflect.ValueOf(&b).Elem()
		err := setBoolVal(v, tc.Value)
		hasErr := err != nil
		if tc.HasErr != hasErr {
			t.Errorf("expected: %t, got: %t", tc.HasErr, hasErr)
		}
		b = v.Bool()
		if b != tc.Bool {
			t.Errorf("expected: %t, got: %t", tc.Bool, b)
		}
	}
}

func equal(s1, s2 interface{}) (bool, error) {
	data1, err := json.Marshal(s1)
	if err != nil {
		return false, err
	}
	data2, err := json.Marshal(s2)
	if err != nil {
		return false, err
	}

	return string(data1) == string(data2), nil
}

func getEnv(key string) (string, error) {
	return os.Getenv(key), nil
}

type mockDva struct{}

func (d *mockDva) getDva(key string) (string, error) {
	return d.Get(key)
}

func (d *mockDva) Start() {}
func (d *mockDva) Stop()  {}
func (d *mockDva) Get(key string) (string, error) {
	switch key {
	case "app.env1":
		return "A", nil
	case "app.env2":
		return "1", nil
	case "app.env3":
		return "true", nil
	case "dummy.yaml":
		return `a: "123"`, nil
	case "dummy.json":
		return `{"a": "123"}`, nil
	default:
		return "", errors.New("dva get error")
	}
}
func (d *mockDva) GetInt(key string) (int, error) {
	return 0, nil
}
func (d *mockDva) GetArr(key string) ([]string, error) {
	return nil, nil
}
func (d *mockDva) GetSecret(key string) (string, error) {
	return "", nil
}
func (d *mockDva) GetStruct(v interface{}) error {
	return nil
}

type Dummy struct {
	A string `json:"a,omitempty" yaml:"a" csv:"a"`
}
