package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const filename = "docsv.md"

func help() string {
	return fmt.Sprintf(
		"HELP\n" +
			"-h help\n" +
			"-r read\n" +
			"-w write [string]\n" +
			"-x done  [int]\n",
	)
}

func title(file *os.File) {
	_, err := file.WriteString("# TODO\n\n")
	if err != nil {
		panic(err)
	}
}

func write(text string) {
	str := read()
	lines := len(str)

	f := file()
	defer f.Close()

	var i int
	if lines == 0 {
		title(f)
		i = lines
	} else {
		i, _ = strconv.Atoi(string(str[lines-1][6]))
	}

	_, err := f.WriteString(fmt.Sprintf("- [ ] %d. %s\n", i+1, text))
	if err != nil {
		panic(err)
	}
}

func read() []string {
	var text []string

	f := file()
	defer f.Close()

	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadString('\n')

		if err != nil {
			break
		}
		const lenPrefix = 8

		if len(line) > lenPrefix {
			text = append(text, line)
		}
	}

	return text
}

func done(x int) {
	text := read()
	f := replaceFile()
	defer f.Close()

	for _, t := range text {
		if strings.Contains(t, strconv.Itoa(x)) {
			t = strings.ReplaceAll(t, "- [ ]", "- [x]")
		}
		f.WriteString(t)
	}
}

func replaceFile() *os.File {
	if err := os.Remove(filename); err != nil {
		panic(err)
	}
	f := file()
	title(f)
	return f
}

func clean() []string {
	removed := []string{}
	text := read()

	f := replaceFile()
	defer f.Close()

	var i int
	for _, t := range text {
		if strings.Contains(t, "- [x]") {
			removed = append(removed, t)
			continue
		}
		i++
		_, err := f.WriteString(fmt.Sprintf("- [ ] %d. %s", i, t[9:]))
		if err != nil {
			panic(err)
		}
	}

	return removed
}

func file() *os.File {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}

func dir() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	dir := home + "/docsv"
	os.Mkdir(dir, 0777)

	if err := os.Chdir(dir); err != nil {
		panic(err)
	}
}

func main() {
	var (
		h bool // help
		r bool // read
		w bool // write
		x int  // done
		c bool // clean
	)
	flag.BoolVar(&h, "h", true, "help")
	flag.BoolVar(&r, "r", false, "read")
	flag.BoolVar(&w, "w", false, "write")
	flag.IntVar(&x, "x", -1, "done")
	flag.BoolVar(&c, "clean", false, "remove all done")
	flag.Parse()

	args := flag.Args()
	dir()

	var result string

	switch {
	case r:
		result = strings.Join(read(), "")

	case w:
		write(strings.Join(args, " "))
		str := read()
		result = str[len(str)-1]

	case x > 0:
		done(x)
		result = strings.Join(read(), "")

	case c:
		result = strings.Join(read(), "")
		result += fmt.Sprintf("CLEANED\n%s", strings.Join(clean(), ""))

	case h:
		result = help()
	default:
	}

	fmt.Print(result)
}
