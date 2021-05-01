package main

import (
	"flag"
)

var (
	flagListen    = flag.String("listen", "[::1]:8080", "Specify what address and port to listen to")
	flagDir       = flag.String("dir", ".", "Directory where to keep files")
	flagNoPersist = flag.Bool("nopersist", false, "Disallow force persist files")
	flagIndex     = flag.String("index", "", "Index file name")
	flagKey       = flag.String("key", "", "Require a key when posting and getting files")
)

type args struct {
	dir             string
	addr            string
	disallowPersist bool
	indexFile       string
	key             string
}

func parseArgs() (args, error) {
	flag.Parse()
	return args{
		dir:             *flagDir,
		addr:            *flagListen,
		disallowPersist: *flagNoPersist,
		indexFile:       *flagIndex,
		key:             *flagKey,
	}, nil
}
