package main

import (
	"flog"
	"fmt"
)

func main() {
	fmt.Println("test")
	flog.Info("wass gud bro")
	flog.SInfo("Im silent!")
	flog.Error("ballin", "big ballin")
	flog.Warn("Nikita", "was", "here")
	flog.Warn("Im also silent! Unless the -dev is a program argument!")

}
