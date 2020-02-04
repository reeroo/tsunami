package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	proxyList []string
	random3     *rand.Rand
	source3     rand.Source
)

func loadProxyList() {
	//Load user agents from file
	file, err := os.Open(*proxyListFile)
	if err != nil {
		//File not found, or whatever, use default UA
		proxyList = append(proxyList, "Tsunami Flooder (https://github.com/ammar/tsunami)")
		fmt.Println(err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			proxyList = append(proxyList, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}
	//Initiate random number generator
	source3 = rand.NewSource(time.Now().UnixNano())
	random3 = rand.New(source3)
}

func getProxyList() string {
	index := int(random.Uint32()) % len(proxyList)
	return proxyList[index]
}