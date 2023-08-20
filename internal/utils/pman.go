package utils

import (
	"github.com/gookit/properties"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	RootPath   = filepath.Dir(b)
)

var PMan *properties.Parser

func NewPman() *properties.Parser {
	pman := properties.NewParser(properties.ParseEnv, properties.ParseInlineSlice)
	fileData, err := os.ReadFile(RootPath + "/../../resources/properties/application.properties")
	if err != nil {
		panic("Properties file does not found")
	}
	err = pman.Parse(string(fileData))
	if err != nil {
		panic(err)
	}
	return pman
}
