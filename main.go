package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

type Check struct {
	timeCheck  int
	statusCode int
	volume     int
}

var Responce_list []Check

func fetchApi(wg *sync.WaitGroup, m *sync.Mutex, api string) {

	m.Lock()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, api, nil)
	if err != nil {
		log.Println(err)
		return
	}

	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	elarsed := time.Since(start).Milliseconds()

	Responce_list = append(Responce_list, Check{
		timeCheck:  int(elarsed),
		statusCode: resp.StatusCode,
		volume:     int(resp.ContentLength),
	})
	m.Unlock()
	wg.Done()

}

func check(api string, amount int) {

	wg := sync.WaitGroup{}
	m := sync.Mutex{}

	for i := 0; i < amount; i++ {
		wg.Add(1)
		go fetchApi(&wg, &m, api)
	}

	wg.Wait()

	log.Print("Результаты чека:\n")
	for i := 0; i < len(Responce_list); i++ {
		fmt.Println("Время запроса: " + strconv.Itoa(Responce_list[i].timeCheck) + "мс Статус код: " + strconv.Itoa(Responce_list[i].statusCode) + " Объем: " + strconv.Itoa(Responce_list[i].volume))
	}
	updateFile()
}

func updateFile() {

	file, err := ioutil.ReadFile("check-results.txt")
	if err != nil {
		log.Println(err)
		return
	}

	fileString := string(file)

	t := time.Now()
	line := "[" + t.Format("2006-11-10 15:04") + "] Результаты чека:\n"
	for _, v := range Responce_list {
		line += "Время запроса: " + strconv.Itoa(v.timeCheck) + "мс Статус код: " + strconv.Itoa(v.statusCode) + " Объем: " + strconv.Itoa(v.volume) + "\n"
	}
	if err := ioutil.WriteFile("check-results.txt", []byte(fileString+line), 0644); err != nil {
		log.Println(err)
		return
	}
	Responce_list = Responce_list[:0]
}
func main() {
	config := NewConfig()

	for {
		check(config.Api, config.Amount_parallelization)
		time.Sleep(time.Duration(config.Repetition_time) * time.Second)
	}
}
