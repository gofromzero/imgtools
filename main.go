package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName("config") // name of config file (without extension)
	viper.SetConfigType("yaml")   // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")
}

type config struct {
	Files []ImgFile `json:"files"`
}

func initConfig(cfg *config) (err error) {
	err = viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("fatal error config file: %w", err)
	}

	err = viper.Unmarshal(cfg) // Find and read the config file
	if err != nil {            // Handle errors reading the config file
		return fmt.Errorf("fatal error config file: %w", err)
	}
	return
}

func main() {
	var err error
	var cfg config
	err = initConfig(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	for i, file := range cfg.Files {
		// check params
		err = file.CheckParams()
		if err != nil {
			log.Fatalf("files %d err: %v", i, err.Error())
		}
	}
	outdir := "out_" + time.Now().Format("2006_01_02_15_04_05")
	CreateDir(outdir)

	total := len(cfg.Files)

	for i, file := range cfg.Files {
		// draw
		err = file.Draw(i, outdir)
		if err != nil {
			log.Printf("%d/%d baseimg[%s] Err: %v!\n", i+1, total, file.Baseimg, err)
		}
		log.Printf("%d/%d Complete!\n", i+1, total)
	}

	log.Println("全部生成结束")

}
