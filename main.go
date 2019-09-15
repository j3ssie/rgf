package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

type pattern struct {
	Flags    string   `json:"flags,omitempty"`
	Pattern  string   `json:"pattern,omitempty"`
	Patterns []string `json:"patterns,omitempty"`
}

type command struct {
	Flags   string
	Pattern string
}

func main() {
	var signs string
	flag.StringVar(&signs, "signs", "", "directory to store signatures")

	// target
	var dir string
	flag.StringVar(&dir, "dir", ".", "directory to search")
	var file string
	flag.StringVar(&file, "file", "", "file to search")

	// adding flag
	var addMode string
	flag.StringVar(&addMode, "add", "", "a new pattern")
	var flags string
	flag.StringVar(&flags, "flags", "-uu -L -C 4 --smart-case", "flag for ripgrep")

	// custom help
	flag.Usage = func() {
		usage()
	}
	checkRg()

	flag.Parse()
	// folder contain signatures
	if signs == "" {
		usr, err := user.Current()
		if err != nil {
			log.Fatal(err)
		}
		signs = usr.HomeDir + "/.rgf/"
	}

	if addMode != "" {
		signName := signs + addMode + ".json"
		// regex of signatures
		regex := flag.Args()[0]
		setSign(signName, flags, regex)
	} else {
		var target string

		if file != "" {
			target = file
		} else {
			if dir == "." {
				dir, err := os.Getwd()
				if err != nil {
					log.Fatal(err)
				}
				// fmt.Println(dir)
				target = dir

			} else {
				target = dir
			}
		}

		specs := flag.Args()
		var spec string
		if len(specs) == 0 {
			spec = "*"
		} else {
			spec = specs[0]
		}

		// list of signatures
		signatures := getSignFiles(signs, spec)
		for _, sign := range signatures {
			doGrep(target, sign)
		}
	}

}

func setSign(signName, flags, regex string) {
	data := pattern{
		Flags:   flags,
		Pattern: regex,
	}
	file, _ := json.MarshalIndent(data, "", " ")
	fmt.Printf("Writing new signatures to: %v \n", signName)
	_ = ioutil.WriteFile(signName, file, 0644)
}

func doGrep(target, signFile string) {
	jsonFile, err := os.Open(signFile)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var pat pattern
	json.Unmarshal(byteValue, &pat)

	if pat.Pattern != "" {
		command := fmt.Sprintf(`rg %s %s %s`, pat.Flags, pat.Pattern, target)
		run(command)
	} else if len(pat.Patterns) > 0 {
		for _, patt := range pat.Patterns {
			command := fmt.Sprintf(`rg %s %s %s`, pat.Flags, patt, target)
			run(command)
		}
	}
}

func getSignFiles(signFolder, spec string) []string {
	os.MkdirAll(signFolder, os.ModePerm)
	realSigns := signFolder + "*.json"
	files, err := filepath.Glob(realSigns)
	if err != nil {
		log.Fatal(err)
	}
	// read each signature file
	var signs []string
	for _, sign := range files {
		if strings.Contains(sign, spec) || spec == "*" {
			signs = append(signs, sign)
		}
	}
	return signs
}

func run(realCmd string) {
	var cmd *exec.Cmd
	command := strings.Split(realCmd, ` `)
	fmt.Println(command)
	cmd = exec.Command(command[0], command[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func checkRg() {
	// _, err := exec.LookPath("/usr/bin/rg")
	cmd := exec.Command("rg", "-h")
	cmdReader, err := cmd.StdoutPipe()
	scanner := bufio.NewScanner(cmdReader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		panic("ripgrep not installed")
		// os.Exit(1)
	}
	go func() {
		for scanner.Scan() {
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ripgrep not installed", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ripgrep not installed", err)
		os.Exit(1)
	}
}

func usage() {
	func() {
		h := "A wrapper around ripgrep to check for various common patterns. \n\n"
		h += "Usage:\n"
		h += "rgf\n"
		h += "rgf -dir /folder/to/grep/\n"
		h += "rgf -file whateverfile\n"
		h += "rgf -dir /folder/to/grep/ url\n\n"
		h += "Add New pattern:\n"
		h += "./rgf.py -add url 'https?://(?:[-\\w.]|(?:%[\\da-fA-F]{2}))+\n"
		fmt.Fprint(os.Stderr, h)
	}()
}
