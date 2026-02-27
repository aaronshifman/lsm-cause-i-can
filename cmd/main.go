package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/aaronshifman/lsm-cause-i-can/pkg/cli"
)

func main() {
	c := cli.NewCli()
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		inter, err := c.Parse(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
			continue
		}

		result, err := c.Execute(inter)
		if errors.Is(err, cli.ErrQuit) {
			return
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}

		fmt.Println(result)
	}
}
