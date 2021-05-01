package main

import (
	"flag"
)

var (
	flagListen = flag.String("listen", "[::1]:8080", "Specify what address and port to listen to")
	flagDir    = flag.String("dir", ".", "Directory where to keep files")
)

type args struct {
	dir  string
	addr string
}

func parseArgs() (args, error) {
	flag.Parse()
	return args{
		dir:  *flagDir,
		addr: *flagListen,
	}, nil
}
