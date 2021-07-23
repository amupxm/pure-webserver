package config

import (
	"encoding/json"
	"log"
	"os"
)

type (
	Config struct {
		Http           httpConfig     `json:"http"`
		DatabaseConfig databaseConfig `json:"database"`
	}
	httpConfig struct {
		Port string `json:"port"`
	}
	databaseConfig struct {
		BucketName string `json:"bucket_name"`
	}
)

var AppConf Config

// TODO : read config from json
func Init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// read config file
	file, err := os.ReadFile(wd + "/config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file, &AppConf)
	if err != nil {
		log.Fatal(err)
	}
	AppConf.DatabaseConfig.BucketName = wd + "/" + AppConf.DatabaseConfig.BucketName
	// hc := httpConfig{
	// 	Port: "8080",
	// }
	// dc := databaseConfig{
	// 	BucketName: wd + "/database.json",
	// }
	// log.Print(wd + "/database.json")
	// AppConf.Http = hc
	// AppConf.DatabaseConfig = dc
}
