package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Config struct {
	Password  string `toml:"password"`
	User      string `toml:"user"`
	Url       string `toml:"url"`
	Nippo_dir string `toml:"nippo_dir"`
}

func main() {
	var config Config
	_, err := toml.DecodeFile("./config.toml", &config)
	if err != nil {
		_, err := toml.DecodeFile("/Users/takedamasaru/.config/blog-cli/config.toml", &config)
		if err != nil {
			panic(err)
		}
	}
	if len(os.Args) < 2 {
		fmt.Println("too few args [post, today, remove]")
	} else if os.Args[1] == "post" {
		post(config)
	} else if os.Args[1] == "today" {
		_, body := getBlog(getPostPath(config.Nippo_dir))
		fmt.Println(body)
	} else if os.Args[1] == "remove" {
		remove(os.Args[2], config)
	}
}
func post(config Config) {
	postpath := getPostPath(config.Nippo_dir)
	title, body := getBlog(postpath)
	values := url.Values{}
	values.Set("body", body)
	values.Add("title", title)
	values.Add("user", config.User)
	values.Add("password", config.Password)
	req, err := http.NewRequest(
		"POST",
		config.Url,
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	defer resp.Body.Close()
}
func getPostPath(nippo_path string) string {
	if len(os.Args) <= 2 {
		files, e := ioutil.ReadDir(nippo_path)
		if e != nil {
			panic(e)
		}
		var lastFile os.FileInfo
		for _, f := range files {
			if lastFile == nil {
				lastFile = f
			} else if lastFile.ModTime().Unix() < f.ModTime().Unix() {
				lastFile = f
			}
		}
		return nippo_path + lastFile.Name()
	} else {
		return os.Args[2]
	}
}
func remove(id string, config Config) {
	values := url.Values{}
	values.Add("id", id)
	values.Add("user", config.User)
	values.Add("password", config.Password)
	req, err := http.NewRequest(
		"DELETE",
		config.Url,
		strings.NewReader(values.Encode()),
	)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	defer resp.Body.Close()
}
func getBlog(filepath string) (string, string) {
	f, er := os.Open(filepath)
	if er != nil {
		fmt.Println("file not found")
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Println("cannot file read")
	}
	body := string(blackfriday.MarkdownCommon(b))
	title := strings.Split(f.Name(), ".")[0]
	if strings.Contains(title, "/") {
		sp := strings.Split(title, "/")
		title = sp[len(sp)-1]
	}
	return title, body
}
