package output

import (
	"fmt"
	"os"
)

func ExitOnErr(err error) {
	if err != nil {
		//nolint:forbidigo
		fmt.Println(err)
	}

	os.Exit(1)
}

func Print(s string) {
	//nolint:forbidigo
	fmt.Println(s)
}
