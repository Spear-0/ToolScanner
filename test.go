package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type arryConfig struct {
	Request struct {
		Method string   `yaml:"method"`
		Path   string   `yaml:"path"`
		Data   []string `yaml:"data"`
	}
	Response struct {
		Pcre_body   []string          `yaml:"pcre_body"`
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

func main() {
	var arrconfig arryConfig
	data, err := ioutil.ReadFile("yaml/test.yaml")
	if err != nil {
		log.Fatalf("%v", err)
	}
	unerr := yaml.Unmarshal(data, &arrconfig)
	if unerr != nil {
		log.Fatalf("%v", unerr)
	}
	log.Println(arrconfig)
	// if arrconfig.Request.Data == nil {
	// 	log.Println("empty")
	// }
	// for k, v := range arrconfig.Request.Data {
	// 	log.Println(k, v)
	// }
	// for k, v := range arrconfig.Response.Pcre_body {
	// 	if v != "" {
	// 		log.Println(k, v)

	// 	}
	// 	log.Println(arrconfig.Response.Pcre_body[k])
	// }
}
