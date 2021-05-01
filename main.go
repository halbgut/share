package main

func main() {
	args, err := parseArgs()
	if err != nil {
		panic(err)
	}
	err = start(args)
	if err != nil {
		panic(err)
	}
}
