package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"

	"gopkg.in/yaml.v2"
)

var VERSION string = "v0.1"

func banner() {
	fmt.Printf("\tToolScanner %s\n", VERSION)
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
	Name    string `yaml:"name"`
	Protect string `yaml:"protect"`
}

func main() {
	banner()
	/**
	* 获取输入参数
	 */
	var host = flag.String("s", "", "Target Address")
	var port = flag.Int("p", 0, "Target Port")
	flag.Parse()

	log.Printf("[*] target server: %s, target port: %d", *host, *port)
	var target = fmt.Sprintf("%s:%d", *host, *port)

	/**
	* 获取yaml文件
	 */
	yamls, err := ioutil.ReadDir("yaml/")
	if err != nil {
		log.Fatal("[-] get yaml fail")
	}

	for _, fileInfo := range yamls {
		if fileInfo.Name()[len(fileInfo.Name())-5:] != ".yaml" {
			return
		}
		var filename = "yaml/" + fileInfo.Name()
		//log.Printf("[*] Using %s", filename)
		data, ioerr := ioutil.ReadFile(filename)
		if ioerr != nil {
			log.Fatal(ioerr)
		}
		var yamlconfig Yaml2Config
		err1 := yaml.Unmarshal(data, &yamlconfig)
		if err1 != nil {
			log.Fatal(err1)
		}
		switch yamlconfig.Protect {
		case "http":
			var url = fmt.Sprintf("http://%s", target)
			if yamlconfig.Request.Path != "" {
				url += yamlconfig.Request.Path
			}
			HTTPParse(yamlconfig, url)
		case "tcp":
			conn, err := net.Dial("tcp", target)
			if err != nil {
				log.Fatalf("[-] %v", err)
			}
			defer conn.Close()
			TCPaser(yamlconfig, conn)
		default:
			log.Fatalf("[-] unknown protect: %s", yamlconfig.Protect)
		}
	}

}

func HTTPParse(yamlconfig Yaml2Config, url string) {
	var isMatchBody = false
	var isMatchStatus = false

	client := &http.Client{}
	req, err := http.NewRequest(yamlconfig.Request.Method, url, nil)
	if err != nil {
		log.Printf("[-] create http client fail")
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[-] get response fail")
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
	// log.Println(resp.Status)
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
func TCPaser(yamlconfig Yaml2Config, conn net.Conn) {
	if yamlconfig.Request.Data != "" {
		_, err := conn.Write([]byte(yamlconfig.Request.Data))
		if err != nil {
			log.Println("Write data fail")
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
	//log.Printf("%s", string(buf[:len]))
	return
}

//Mf5pknzyxgP1YTmhG
