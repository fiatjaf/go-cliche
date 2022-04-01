package main

import (
	"log"
	"time"

	"github.com/fiatjaf/go-cliche"
)

func main() {
	c := &cliche.Control{
		DataDir: ".",
		JARPath: "/home/fiatjaf/comp/cliche/target/scala-2.13/cliche-assembly-0.1.0.jar",

		LogStderr:           true,
		LogIrrelevantStdout: true,
	}

	log.Print(c.Start())
	log.Print(c.GetInfo())

	time.Sleep(5 * time.Second)
	log.Print(c.GetInfo())
}
