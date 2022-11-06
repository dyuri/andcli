package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

const app = "andcli"

var (
	cfgDir  = "."
	cfgFile = "config.yaml"

	// build vars
	commit = ""
	date   = ""
	arch   = ""
	gover  = ""
	tag    = ""
)

func init() {
	initConfig()
	log.SetFlags(0)
}

func main() {
	var vaultFile, vaultType string
	var showVersion bool

	flag.StringVar(&vaultFile, "f", "", "Path to the encrypted vault")
	flag.StringVar(&vaultType, "t", "", "Vault type (andotp, aegis)")
	flag.BoolVar(&showVersion, "v", false, "Show current version")
	flag.Parse()

	if showVersion {
		fmt.Printf(
			"%s %s %s (%s) built on %s with go version %s\n",
			app, tag, arch, commit, date, gover,
		)
		os.Exit(0)
	}

	c := configFromFile(vaultFile, vaultType)
	if c.File == "" {
		log.Fatal("[ERR] missing input file, specify one with -f")
	}

	if c.Type == "" {
		log.Fatal("[ERR] missing vault type, specify one with -t")
	}

	abs, err := filepath.Abs(c.File)
	if err != nil {
		log.Fatal("[ERR] ", err)
	}

	entries, err := decrypt(abs, c.Type)
	if err != nil {
		log.Fatal("[ERR] ", err)
	}

	if err := c.save(); err != nil {
		log.Fatal("[ERR] ", err)
	}

	termenv.ClearScreen()

	p := tea.NewProgram(newModel(c.File, entries))
	if err := p.Start(); err != nil {
		log.Fatal("[ERR] ", err)
	}
}

func initConfig() {
	var err error

	cfgDir, err = os.UserConfigDir()
	if err != nil {
		log.Fatal("[ERR] open config dir: ", err)
	}

	cfgDir = filepath.Join(cfgDir, app)
	cfgFile = filepath.Join(cfgDir, cfgFile)
}
