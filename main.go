package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден")
	}
}

type Config struct {
	Api                    string
	Repetition_time        int
	Amount_parallelization int
}

func NewConfig() *Config {

	api, exists := os.LookupEnv("API")
	if exists {
		log.Println("API: " + api)
	}

	repe_time, exists := os.LookupEnv("REPETITION_TIME")
	if exists {
		log.Println("REPETITION_TIME: " + repe_time)
	}

	repetition_time, err := strconv.Atoi(repe_time)
	if err != nil {
		panic(err)
	}

	amount_para, exists := os.LookupEnv("AMOUNT_PARALLELIZATION")
	if exists {
		log.Println("AMOUNT_PARALLELIZATION: " + amount_para)
	}

	amount_parallelization, err := strconv.Atoi(amount_para)
	if err != nil {
		panic(err)
	}

	return &Config{
		Api:                    api,
		Repetition_time:        repetition_time,
		Amount_parallelization: amount_parallelization,
	}
}

var Responce_list []int

func fetchApi(req *http.Request) {

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	Responce_list = append(Responce_list, resp.StatusCode)
}

func main() {
	config := NewConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for {
		for i := 0; i < config.Amount_parallelization; i++ {

			req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.Api, nil)
			if err != nil {
				log.Println(err)
			}

			go fetchApi(req)
		}

		log.Println("Результаты тестрования:")
		fmt.Println(Responce_list)
		time.Sleep(time.Duration(config.Repetition_time) * time.Second)
	}
}
