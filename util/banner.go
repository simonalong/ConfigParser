package util

import (
	"fmt"
	"os"
)

var Banner = DefaultBanner

func PrintBanner() {
	fmt.Printf("%s\n", Banner)
}

func LoadBanner(filePath string) {
	if b, err := os.ReadFile(filePath); err == nil {
		str := string(b)
		fmt.Printf("%s\n", str)
	} else {
		fmt.Printf("%s\n", Banner)
	}
}

var DefaultBanner = "" +
	"                 ██████╗    ██████╗   ██╗       ███████╗\n" +
	"                ██╔════╝   ██╔═══██╗  ██║       ██╔════╝\n" +
	"                ██║  ███╗  ██║   ██║  ██║       █████╗  \n" +
	"                ██║   ██║  ██║   ██║  ██║       ██╔══╝  \n" +
	"                ╚██████╔╝  ╚██████╔╝  ███████╗  ███████╗\n" +
	"                 ╚═════╝    ╚═════╝   ╚══════╝  ╚══════╝\n                                  "
