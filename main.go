package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"sync"

	"gopkg.in/yaml.v2"
)

var VERSION string = "v0.1"

func banner() {
	fmt.Printf("ToolScanner %s\n", VERSION)
}

type Yaml2Config struct {
	Request struct {
		Method string `yaml:"method"`
		Path   string `yaml:"path"`
		Data   string `yaml:"data"`
	}
	Response struct {
		Pcre_body   string            `yaml:"pcre_body"`
		Pcre_status string            `yaml:"pcre_status"`
		Pcre_header map[string]string `yaml:"pcre_header"`
	}
	Tool struct {
		Tool_name    string `yaml:"tool_name"`
		Tool_version string `yaml:"tool_version"`
	}
	Name     string `yaml:"name"`
	Protocol string `yaml:"protocol"`
}

type ConfigQueue struct {
	items []Yaml2Config
	lock  sync.RWMutex
}

func (q *ConfigQueue) CreateQueue() *ConfigQueue {
	q.items = []Yaml2Config{}
	return q
}

func (q *ConfigQueue) Push(item Yaml2Config) {
	q.lock.Lock()
	q.items = append(q.items, item)
	q.lock.Unlock()
}

func (q *ConfigQueue) Pop() *Yaml2Config {
	q.lock.Lock()
	item := q.items[len(q.items)-1]
	q.items = q.items[:len(q.items)-1]
	q.lock.Unlock()
	return &item
}

func (q *ConfigQueue) Size() int {
	return len(q.items)
}

func (q *ConfigQueue) IsEmpty() bool {
	return q.Size() == 0
}

func main() {

	banner()

	var host = flag.String("s", "", "Target Address")
	var port = flag.Int("p", 0, "Target Port")
	flag.Parse()

	log.Println("[*] Init config queue...")
	var http ConfigQueue
	var httpQueue = http.CreateQueue()
	var tcp ConfigQueue
	var tcpQueue = tcp.CreateQueue()

	log.Printf("[*] target server: %s, target port: %d", *host, *port)
	var target = fmt.Sprintf("%s:%d", *host, *port)

	log.Println("[*] get yaml file")
	yamls, err := ioutil.ReadDir("yaml/")
	if err != nil {
		log.Fatal("[-] get yaml file fail")
	}

	for _, fileInfo := range yamls {
		if fileInfo.Name()[len(fileInfo.Name())-5:] != ".yaml" {
			return
		}
		var filename = "yaml/" + fileInfo.Name()
		data, ioerr := ioutil.ReadFile(filename)
		if ioerr != nil {
			log.Fatal(ioerr)
		}
		var yamlconfig Yaml2Config
		err1 := yaml.Unmarshal(data, &yamlconfig)
		if err1 != nil {
			log.Fatal(err1)
		}
		switch yamlconfig.Protocol {
		case "http":
			httpQueue.Push(yamlconfig)
		case "tcp":
			tcpQueue.Push(yamlconfig)
		default:
			log.Printf("[-] Currently, the current protocol is not supported[%s]", yamlconfig.Name)
		}

	}

	log.Println("[*] start connect target server...")
	if !httpQueue.IsEmpty() {
		ExecuteHTTPQueue(httpQueue, target)
	}
	if !tcpQueue.IsEmpty() {
		ExecuteTCPQueue(tcpQueue, target)
	}
}
func ExecuteHTTPQueue(q *ConfigQueue, target string) {
	i := q.Size()
	for {
		item := q.Pop()
		HTTPParse(*item, "http://"+target+item.Request.Path)
		i -= 1
		if i == 0 {
			break
		}
	}
}

func ExecuteTCPQueue(q *ConfigQueue, target string) {
	i := q.Size()
	for {
		item := q.Pop()
		TCPPaser(*item, target)
		i -= 1
		if i == 0 {
			break
		}
	}
}
func HTTPParse(yamlconfig Yaml2Config, url string) {
	var isMatchBody = false
	var isMatchStatus = false

	client := &http.Client{}
	req, err := http.NewRequest(yamlconfig.Request.Method, url, nil)
	if err != nil {
		log.Printf("[-] %v", err)
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[-] %v", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if yamlconfig.Response.Pcre_body != "" {
		match, err := regexp.Match(yamlconfig.Response.Pcre_body, body)
		if err != nil {
			log.Printf("[-] %v", err)
			return
		}
		if match == true {
			isMatchBody = true
		}
	} else {
		isMatchBody = true
	}
	if yamlconfig.Response.Pcre_status != "" {
		match, err := regexp.MatchString(yamlconfig.Response.Pcre_status, resp.Status)
		if err != nil {
			log.Fatalf("[-] %v", err)
			return
		}
		if match {
			isMatchStatus = true
		}
	} else {
		isMatchStatus = true
	}
	if len(yamlconfig.Response.Pcre_header) != 0 {
		for pcreKey, pcereVal := range yamlconfig.Response.Pcre_header {
			_, ok := resp.Header[pcreKey]
			if ok {
				match, err := regexp.MatchString(pcereVal, resp.Header[pcreKey][0])
				if err != nil {
					log.Printf("[-] %v", err)
					return
				}
				if !match {
					return
				}
			} else {
				return
			}
		}
		if isMatchBody && isMatchStatus {
			log.Printf("[%s %s]", yamlconfig.Tool.Tool_name, yamlconfig.Tool.Tool_version)
		}
	} else {
		log.Printf("[%s %s]", yamlconfig.Tool.Tool_name, yamlconfig.Tool.Tool_version)
	}
}
func TCPPaser(yamlconfig Yaml2Config, target string) {
	conn, err := net.Dial("tcp", target)
	if err != nil {
		log.Fatalf("[-] %v", err)
	}
	defer conn.Close()
	if yamlconfig.Request.Data != "" {
		_, err := conn.Write([]byte(yamlconfig.Request.Data))
		if err != nil {
			log.Printf("[-] %v", err)
			return
		}
	}
	var buf [2048]byte
	len, err := conn.Read(buf[:])
	conn.Close()
	if err != nil {
		log.Println("Read data fail")
		return
	}
	match, err := regexp.Match(yamlconfig.Response.Pcre_body, buf[:len])
	if match {
		log.Printf("[%s %s]", yamlconfig.Tool.Tool_name, yamlconfig.Tool.Tool_version)
	}
	return
}
func UDPPaser(yamlconfig Yaml2Config, target string) {}

func DNSPaser(yamlconfig Yaml2Config, target string) {}