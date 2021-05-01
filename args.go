package main

import (
	"flag"
)

var (
	flagListen    = flag.String("listen", "[::1]:8080", "Specify what address and port to listen to")
	flagDir       = flag.String("dir", ".", "Directory where to keep files")
	flagNoPersist = flag.Bool("nopersist", false, "Disallow force persist files")
)

type args struct {
	dir             string
	addr            string
	disallowPersist bool
}

func parseArgs() (args, error) {
	flag.Parse()
	return args{
		dir:             *flagDir,
		addr:            *flagListen,
		disallowPersist: *flagNoPersist,
	}, nil
}
