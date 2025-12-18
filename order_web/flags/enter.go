package flags

import (
	"flag"
)

type Options struct {
	File    string
	Version bool
}

var FileOption = new(Options)

func Parse() {
	flag.StringVar(&FileOption.File, "f", "settings.yaml", "配置文件")
	flag.BoolVar(&FileOption.Version, "v", false, "版本")
	flag.Parse()
}
