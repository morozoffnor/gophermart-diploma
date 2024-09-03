package config

import (
	"flag"
	"os"
)

type Config struct {
	Addr              string
	DatabaseURI       string
	AccrualSystemAddr string
}

func (c *Config) UpdateByFlags(f *Flags) {
	flag.Parse()
	c.Addr = f.Addr
	c.DatabaseURI = f.DatabaseURI
	c.AccrualSystemAddr = f.AccrualSystemAddr
}

func (c *Config) UpdateByEnv() {
	a := os.Getenv("RUN_ADDRESS")
	if a != "" {
		c.Addr = a
	}
	d := os.Getenv("DATABASE_URI")
	if d != "" {
		c.DatabaseURI = d
	}
	asa := os.Getenv("ACCRUAL_SYSTEM_ADDRESS")
	if asa != "" {
		c.AccrualSystemAddr = asa
	}
}

func New() *Config {
	c := &Config{}
	c.UpdateByFlags(flags)
	c.UpdateByEnv()
	return c
}
