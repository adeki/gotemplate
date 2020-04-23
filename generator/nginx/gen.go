package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

type Server struct {
	ServerName string `yaml:"server_name"`
	Host       string `yaml:"host"`
	Port       string `yaml:"port"`
}

type Servers map[string]Server

type Environments map[string]Servers

type Config struct {
	TemplateDir  string       `yaml:"template_dir"`
	OutputDir    string       `yaml:"output_dir"`
	Environments Environments `yaml:"environments"`
}

func (c Config) OutputEnvDir(env string) string {
	return c.OutputDir + "/" + env + "/nginx/conf.d"
}

var ConfigFile string

func init() {
	flag.StringVar(&ConfigFile, "c", "ngxconf.yml", "Configuration file to use.")
	flag.StringVar(&ConfigFile, "config", "ngxconf.yml", "Configuration file to use.")
}

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatalln(err)
	}
	var config Config
	if err := yaml.Unmarshal(b, &config); err != nil {
		log.Fatalln(err)
	}

	for envName, servers := range config.Environments {
		for serverName, server := range servers {
			tmplFileName := "nginx." + serverName + ".tmpl"
			t := template.Must(template.ParseFiles(config.TemplateDir + "/" + tmplFileName))
			output, err := openFile(config.OutputEnvDir(envName), serverName)
			if err != nil {
				log.Println(err)
				return
			}
			if err := t.Execute(output, server); err != nil {
				log.Println(err)
				return
			}
			output.Close()
		}
	}
}

func openFile(outputDir, server string) (*os.File, error) {
	if f, err := os.Stat(outputDir); os.IsNotExist(err) || !f.IsDir() {
		if err := os.MkdirAll(outputDir, 0777); err != nil {
			return nil, err
		}
	}
	return os.OpenFile(outputDir+"/"+server+".conf", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
}
