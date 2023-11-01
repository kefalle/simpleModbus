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

func defaultUint16Action(val interface{}, t *controller.Tag) {
	if t.LastValue != val {
		v := val.(uint16)
		log.Printf("%s = %d", t.Name, v)
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

	// Эти регистры как в контроллере +1
	ctrl.AddTag(&controller.Tag{Name: "temp_floor", Address: 513, Method: controller.READ_FLOAT, Action: defaultFloat32Action})
	ctrl.AddTag(&controller.Tag{Name: "temp_otopl", Address: 515, Method: controller.READ_FLOAT, Action: defaultFloat32Action})
	ctrl.AddTag(&controller.Tag{Name: "temp_boiler", Address: 517, Method: controller.READ_FLOAT, Action: defaultFloat32Action})
	ctrl.AddTag(&controller.Tag{Name: "temp_inout", Address: 519, Method: controller.READ_FLOAT, Action: defaultFloat32Action})

	// Регистры как в контроллере
	ctrl.AddTag(&controller.Tag{Name: "status", Address: 520, Method: controller.READ_UINT, Action: defaultUint16Action})
	ctrl.AddTag(&controller.Tag{Name: "servo_otopl", Address: 521, Method: controller.READ_UINT, Action: defaultUint16Action})
	ctrl.AddTag(&controller.Tag{Name: "servo_floor", Address: 550, Method: controller.READ_UINT, Action: defaultUint16Action})

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
