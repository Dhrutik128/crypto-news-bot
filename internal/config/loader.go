package config

import (
	"encoding/json"
	"flag"
	"fmt"
	converter "github.com/ghodss/yaml"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type Configuration interface {
	Load()
	Path() string
	Valid() error
}

func readFile(cfgFile string) ([]byte, error) {
	// read configuration
	return ioutil.ReadFile(cfgFile)
}

// convert yaml into json
func parseConfigurationJson(yamlData []byte) ([]byte, error) {
	jsonDoc, err := converter.YAMLToJSON(yamlData)
	if err != nil {
		fmt.Printf("Error converting YAML to JSON: %s\n", err.Error())
		return nil, err
	}
	return jsonDoc, nil
}

// retrieves configuration in json format (converted from yaml into json)
func Configure(conf Configuration) {
	// read the configuration file
	cfgFile := conf.Path()
	data, err := readFile(cfgFile)
	if err != nil {
		panic(err)
	}
	// check if file is in yaml format
	if cfgFile[len(cfgFile)-4:] == "yaml" {
		// convert the yaml data into json data
		data, err = parseConfigurationJson(data)
	}
	// unmarshal that json data into the concrete configuration struct
	unmarshalJson(data, conf)
	if Bool("DEBUG", false) {
		fmt.Printf("%+v\n", conf)
	}
	if err := conf.Valid(); err != nil {
		log.Fatal(err)
	}
}

// unmarshal in into conf
func unmarshalJson(in []byte, conf Configuration) {
	if err := json.Unmarshal(in, &conf); err != nil {
		log.Panic(err)
	}
}

// reads configflag from env (using fallback)
// creates identifier '--CONFIGFLAG='
// checks os.Args for flag occurrence.
// finally sets cfgfile to the parsed argument value.
func FileNameFromFlag(envFlagFallBack string, defaultFileName string, flagUsage string) string {
	var cfgfile string
	configFlag := String("CONFIGFLAG", envFlagFallBack)
	configFlagIdentifier := fmt.Sprintf("--%s=", configFlag)
	flag.StringVar(&cfgfile, configFlag, defaultFileName, flagUsage)
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, configFlagIdentifier) {
			err := flag.CommandLine.Parse([]string{arg})
			if err != nil {
				panic(err)
			}
			break
		}
	}
	return cfgfile
}
