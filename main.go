package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var NilFile = errors.New("file is nil")

func main() {
	core := core{}

	if len(os.Args) == 1 {
		return
	}

	core.load()

	switch os.Args[1] {
	case "add":
		core.write(os.Args[2])

	case "read":
		text, err := core.read()
		if err != nil {
			fmt.Println(err)
		}

		for _, slot := range text {
			fmt.Println(slot)
		}

	default:
		fmt.Fprintf(os.Stdout, "%s no comand \n", os.Args[1])
	}
	core.save()
}

type core struct {
	id   int
	data []string
}

func (c core) read() (buf []string, err error) {
	if len(c.data) == 0 {
		return []string{}, NilFile
	}

	return c.data, nil
}

func (c *core) write(text string) (err error) {
	if text == "" {
		return NilFile
	}

	c.data = append(c.data, text)
	return nil
}

func (c core) save() (err error) {
	if len(c.data) == 0 {
		return NilFile
	}

	var str string
	for _, i := range c.data {
		str += i + "\n"
	}

	err = os.WriteFile("example.txt", []byte(str), 0666)
	return nil
}

func (c *core) load() (err error) {
	/*if len(c.data) == 0 {
		return NilFile
	}*/

	data, err := os.ReadFile("example.txt")

	parts := strings.Split(string(data), ",")

	for _, i := range parts {
		c.data = append(c.data, i)
	}
	return nil
}
