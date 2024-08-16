package environment

import (
	"os"
	"testing"
)

type Config struct {
	Value string `json:"value"`
}

type Environment struct {
	StringValue        string  `env:"STRING_VALUE"`
	IntValue           int     `env:"INT_VALUE"`
	Int64Value         int64   `env:"INT64_VALUE"`
	Int32Value         int32   `env:"INT32_VALUE"`
	Int16Value         int16   `env:"INT16_VALUE"`
	Int8Value          int8    `env:"INT8_VALUE"`
	UintValue          uint    `env:"UINT_VALUE"`
	Uint64Value        uint64  `env:"UINT64_VALUE"`
	Uint32Value        uint32  `env:"UINT32_VALUE"`
	Uint16Value        uint16  `env:"UINT16_VALUE"`
	Uint8Value         uint8   `env:"UINT8_VALUE"`
	Float32Value       float32 `env:"FLOAT32_VALUE"`
	Float64Value       float64 `env:"FLOAT64_VALUE"`
	BoolValue          bool    `env:"BOOL_VALUE"`
	PointerValue       *Config `env:"POINTER_VALUE"`
	PointerScalarValue *int    `env:"POINTER_SCALAR_VALUE"`
}

func TestReadEnvironment(t *testing.T) {
	_ = os.Setenv("STRING_VALUE", "string")
	_ = os.Setenv("INT_VALUE", "1")
	_ = os.Setenv("INT64_VALUE", "2")
	_ = os.Setenv("INT32_VALUE", "3")
	_ = os.Setenv("INT16_VALUE", "4")
	_ = os.Setenv("INT8_VALUE", "5")
	_ = os.Setenv("UINT_VALUE", "6")
	_ = os.Setenv("UINT64_VALUE", "7")
	_ = os.Setenv("UINT32_VALUE", "8")
	_ = os.Setenv("UINT16_VALUE", "9")
	_ = os.Setenv("UINT8_VALUE", "10")
	_ = os.Setenv("FLOAT32_VALUE", "1.1")
	_ = os.Setenv("FLOAT64_VALUE", "1.2")
	_ = os.Setenv("BOOL_VALUE", "true")
	_ = os.Setenv("POINTER_VALUE", `{"value":"value"}`)
	_ = os.Setenv("POINTER_SCALAR_VALUE", "11")

	env := Environment{PointerValue: &Config{}, PointerScalarValue: new(int)}

	err := Read(&env)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if env.StringValue != "string" {
		t.Errorf("expected %v, got %v", "string", env.StringValue)
	}

	if env.IntValue != 1 {
		t.Errorf("expected %v, got %v", 1, env.IntValue)
	}

	if env.Int64Value != 2 {
		t.Errorf("expected %v, got %v", 2, env.Int64Value)
	}

	if env.Int32Value != 3 {
		t.Errorf("expected %v, got %v", 3, env.Int32Value)
	}

	if env.Int16Value != 4 {
		t.Errorf("expected %v, got %v", 4, env.Int16Value)
	}

	if env.Int8Value != 5 {
		t.Errorf("expected %v, got %v", 5, env.Int8Value)
	}

	if env.UintValue != 6 {
		t.Errorf("expected %v, got %v", 6, env.UintValue)
	}

	if env.Uint64Value != 7 {
		t.Errorf("expected %v, got %v", 7, env.Uint64Value)
	}

	if env.Uint32Value != 8 {
		t.Errorf("expected %v, got %v", 8, env.Uint32Value)
	}

	if env.Uint16Value != 9 {
		t.Errorf("expected %v, got %v", 9, env.Uint16Value)
	}

	if env.Uint8Value != 10 {
		t.Errorf("expected %v, got %v", 10, env.Uint8Value)
	}

	if env.Float32Value != 1.1 {
		t.Errorf("expected %v, got %v", 1.1, env.Float32Value)
	}

	if env.Float64Value != 1.2 {
		t.Errorf("expected %v, got %v", 1.2, env.Float64Value)
	}

	if env.BoolValue != true {
		t.Errorf("expected %v, got %v", true, env.BoolValue)
	}

	if env.PointerValue.Value != "value" {
		t.Errorf("expected %v, got %v", "value", env.PointerValue.Value)
	}

	if *env.PointerScalarValue != 11 {
		t.Errorf("expected %v, got %v", 11, *env.PointerScalarValue)
	}

	os.Clearenv()
}

type EnvironmentDefault struct {
	StringValue string `env:"STRING_VALUE" default:"default"`
}

func TestReadEnvironmentDefault(t *testing.T) {
	env := EnvironmentDefault{}

	err := Read(&env)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if env.StringValue != "default" {
		t.Errorf("expected %v, got %v", "", env.StringValue)
	}
}

func TestReadEnvironmentDefaultEmpty(t *testing.T) {
	_ = os.Setenv("STRING_VALUE", "")

	env := EnvironmentDefault{}

	err := Read(&env)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if env.StringValue != "default" {
		t.Errorf("expected %v, got %v", "default", env.StringValue)
	}
}
