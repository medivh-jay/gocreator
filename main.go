// starter 项目支持工具, 修改包名和引用名
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var clone = "https://github.com/medivh-jay/starter.git"

func main() {
	var module string
	var help bool
	flag.StringVar(&module, "m", "", "type your custom module name")
	flag.BoolVar(&help, "h", false, "get help")
	flag.Parse()

	if help || module == "" {
		flag.Usage()
		return
	}
	fmt.Println(module)
	cmd := exec.Command("git", "clone", clone, module)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	_ = cmd.Run()

	path, _ := filepath.Abs(fmt.Sprintf("%s/%s/", filepath.Dir(os.Args[0]), module))
	var confirm string
	fmt.Printf("the directory will be traversed: %s, and all files ending in .go will be modified, Confirm?[Y/N]: ", path)
	_, _ = fmt.Scanln(&confirm)
	if strings.ToUpper(confirm) != "Y" {
		return
	}

        git, _ := filepath.Abs(path + "/.git")
	_ = os.RemoveAll(git)
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			buf = bytes.ReplaceAll(buf, []byte(`"starter/`), []byte(fmt.Sprintf(`"%s/`, module)))
			err = ioutil.WriteFile(path, buf, 0755)
			if err != nil {
				return err
			}

			_ = exec.Command("go", "fmt", path).Run()
			fmt.Println(path, " ", " complete!")
		}

		if info.Name() == "go.mod" {
			buf, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			content := string(buf)
			content = strings.ReplaceAll(content, `starter`, fmt.Sprintf(`%s`, module))
			err = ioutil.WriteFile(path, []byte(content), 0755)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		log.Println(err)
	}

}
