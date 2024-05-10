package structconf_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kevinfalting/structconf"
)

func ExampleParse() {
	os.Setenv("NAME", "Vikki")

	type Config struct {
		Name string `conf:"env:NAME"`
	}

	var cfg Config

	if err := structconf.Parse(context.Background(), &cfg); err != nil {
		log.Println("failed to Parse:", err)
	}

	fmt.Printf("%+v\n", cfg)

	// Output:
	// {Name:Vikki}
}
