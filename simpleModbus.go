package main

import (
	"log"
	"net/http"
	"os"
	"simpleModbus/controller"
)

const APP = "simpleModbus"
const VERSION = "0.0.1"

func defaultFloat32Action(val interface{}, t *controller.Tag) {
	if t.LastValue != val {
		v := val.(float32)
		log.Printf("%s = %f", t.Name, v)
		t.LastValue = val
	}
}

// Инициализация модбас контроллера
func initController() (ctrl *controller.Controller, err error) {
	ctrl, err = controller.New(&controller.Configuration{
		Url: "rtuovertcp://192.168.1.200:8899",
	})
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	ctrl.AddTag(&controller.Tag{
		Name:    "temp_1",
		Address: 513,
		Method:  controller.READ_FLOAT,
		Action:  defaultFloat32Action,
	})

	ctrl.AddTag(&controller.Tag{
		Name:    "temp_2",
		Address: 515,
		Method:  controller.READ_FLOAT,
		Action:  defaultFloat32Action,
	})

	return
}

// Инициализация сервера http для выдачи состояния и метрик
func initHttpServer(ctrl *controller.Controller) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/tags", TagsHahdler(ctrl))
	mux.Handle("/metrics", MetricsHandler())

	return mux
}

func main() {
	log.Println("Starting...")

	// Инициализация модбас конроллера
	ctrl, err := initController()
	if err != nil {
		log.Println("Can not listen http")
		os.Exit(1)
	}

	// Запуск полера
	go ctrl.Poll()
	defer ctrl.Close()

	// Инициализация сервера
	mux := initHttpServer(ctrl)
	log.Println("Listening...")
	err = http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Println("Can not listen http")
		os.Exit(1)
	}

	os.Exit(0)
}
