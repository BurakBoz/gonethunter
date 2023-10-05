/*
Copyright (c) 2023 Burak Boz git@burakboz.net

This program is created for educational purposes only. The developer is not
responsible for its usage. Scanning any network without permission may be
considered illegal in your country. Users are solely responsible for their actions.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
    "bufio"
    "crypto/tls"
    "flag"
    "fmt"
    "io"
    "net"
    "net/http"
    "os"
    "os/exec"
    "strings"
    "sync"
    "time"
)

var (
	mu              sync.Mutex
	progressCounter int
	foundExec       *string
	outputFileName  *string
	appendFile      *bool
	silentMode      *bool
	stopOnFound     *bool
)

func showProgress(total, current int, startTime time.Time, requestTimeout time.Duration, remainingIPs int) {
    if *silentMode {
        return;
    }
	mu.Lock()
	defer mu.Unlock()

	progressCounter++

	if progressCounter == total {
		fmt.Printf("Progress: %d/%d - Elapsed Time: %s\n", current, total, time.Since(startTime).Round(time.Second))
	} else if progressCounter < total {
		elapsed := time.Since(startTime)
		if current > 0 {
			avgProcessingTime := elapsed / time.Duration(current)
			remainingTime := avgProcessingTime * time.Duration(remainingIPs)

			fmt.Printf("Progress: %d/%d - Elapsed Time: %s - Estimated Remaining Time: %s\n", current, total, elapsed.Round(time.Second), remainingTime.Round(time.Second))
		} else {
			fmt.Printf("Progress: %d/%d - Elapsed Time: %s\n", current, total, elapsed.Round(time.Second))
		}
	}
}

func main() {
    inputFileName := flag.String("input", "iplist.txt", "Input file name (ip list)")
    outputFileName = flag.String("output", "found.txt", "Input file name (ip list) example: /etc/hosts")
    hostname := flag.String("hostname", "google.com", "Hostname to check")
    searchText := flag.String("search", "www.google.com", "Search text in content")
    foundExec = flag.String("foundExec", "", "Execute command when ip address found. [domain, ip] example: 'echo %s %s >/dev/null'")
    requestTimeout := flag.Duration("timeout", 10*time.Second, "Request timeout per each request")
    concurrentThreads := flag.Int("threads", 2000, "Concurrent threads")
    appendFile = flag.Bool("appendFile", false, "Append to results")
    silentMode = flag.Bool("silentMode", false, "Silent mode")
    stopOnFound = flag.Bool("stopOnFound", false, "Stop on found")
    flag.Parse()

    if !*silentMode {
        fmt.Println("# - - - - - - - - - - - - - - - - #")
        fmt.Println("# - NetHunter IP Scanner v0.0.1 - #")
        fmt.Println("# - - - - - - - - - - - - - - - - #")
        fmt.Println("")
        fmt.Println("Copyright (c) 2023 Burak Boz git@burakboz.net");
        fmt.Println("This program is created for educational purposes only. The developer is not");
        fmt.Println("responsible for its usage. Scanning any network without permission may be");
        fmt.Println("considered illegal in your country. Users are solely responsible for their actions.");
        fmt.Println("")
    }

	current := 0

	ipAddresses, err := readLines(*inputFileName)
	if err != nil {
		if !*silentMode {
		    fmt.Println("File read error:", err)
		}
		return
	}
	totalIPs := len(ipAddresses)

	var wg sync.WaitGroup
	var matchedIPs []string

	startTime := time.Now()

	queue := make(chan string, totalIPs)

	for _, ip := range ipAddresses {
		queue <- ip
	}

	close(queue)

	progressCh := make(chan struct{}, *concurrentThreads)

	for i := 0; i < *concurrentThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for ip := range queue {
				if checkHost(ip, *hostname, *searchText, *requestTimeout) {
					if !*silentMode {
					    fmt.Printf("IP address %s matches (proto HTTPS)\n", ip)
					}
					matchedIPs = append(matchedIPs, ip)
				}

				current += 1;
				remaining := totalIPs-current;

				progressCh <- struct{}{}
				showProgress(totalIPs, current, startTime, *requestTimeout, remaining)
			}
		}()
	}

	go func() {
		wg.Wait()
		close(progressCh)
	}()

	for range progressCh {
		// progressCh
	}

}

func readLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func checkHost(ip, hostname, searchText string, requestTimeout time.Duration) bool {
	target := net.JoinHostPort(ip, "443")

	httpsTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpsClient := &http.Client{
		Timeout:   requestTimeout,
		Transport: httpsTransport,
	}
	httpsReq, err := http.NewRequest("GET", "https://"+target, nil)
	if err != nil {
		return false
	}

	httpsResp, err := httpsClient.Do(httpsReq)

	if err != nil {
		if strings.Contains(err.Error(), "context deadline exceeded") {
			return false
		}
		return false
	}

	defer httpsResp.Body.Close()

	body, err := io.ReadAll(httpsResp.Body)
	if err != nil {
		return false
	}

	if searchText != "" && !strings.Contains(string(body), searchText) {
		return false
	}

	if !*silentMode {
	    fmt.Printf("Found IP: %s\n\n", ip)
	}
	if *appendFile {
	    if err := appendToFile(*outputFileName, hostname, []string{ip}); err != nil {
        	if !*silentMode {
        	    fmt.Println("File append error:", err)
        	}
        }
	} else {
	    if err := writeToFile(*outputFileName, hostname, []string{ip}); err != nil {
            if !*silentMode {
                fmt.Println("File write error:", err)
            }
        }
	}

    if *foundExec != "" {
        cmd := exec.Command("sh", "-c", fmt.Sprintf(*foundExec, hostname, ip))
        if err := cmd.Run(); err != nil {
        	if !*silentMode {
        	    fmt.Println("Shell command execution error:", err)
            }
        }
    }

    if *stopOnFound {
        os.Exit(1)
    }

	return true
}

func writeToFile(filename, hostname string, ips []string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		_, err := writer.WriteString(fmt.Sprintf("%s:%s\n", hostname, ip))
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}

func appendToFile(filename, hostname string, ips []string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, ip := range ips {
		_, err := writer.WriteString(fmt.Sprintf("%s    %s\n", ip, hostname))
		if err != nil {
			return err
		}
	}
	writer.Flush()
	return nil
}
