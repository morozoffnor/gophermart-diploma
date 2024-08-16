package config

import "flag"

type Flags struct {
	Addr              string
	DatabaseURI       string
	AccrualSystemAddr string
}

var flags = newFlags()

func (f *Flags) initFlags() {
	flag.StringVar(&f.Addr, "a", "localhost:8080", "server address")
	flag.StringVar(&f.DatabaseURI, "d", "", "database uri")
	flag.StringVar(&f.AccrualSystemAddr, "r", "", "accrual system address")
}

func newFlags() *Flags {
	f := &Flags{}
	f.initFlags()
	return f
}
