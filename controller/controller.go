package controller

import (
	"github.com/VictoriaMetrics/metrics"
	"github.com/mcuadros/go-defaults"
	"github.com/simonvetter/modbus"
	"log"
	"net"
	"sync"
	"time"
)

type OperationType uint

const (
	READ_UINT  = iota
	READ_FLOAT = iota
	//WRITE_UINT = iota
	//WRITE_UINT = iota
)

type logger struct {
	prefix       string
	customLogger *log.Logger
}

type Configuration struct {
	DeviceId    uint8 `default:"16"`
	Url         string
	Speed       uint          `default:"19200"`
	Timeout     time.Duration `default:"1s"`
	PollingTime time.Duration `default:"1s"`
	ReadPeriod  time.Duration `default:"10ms"`
	// Logger provides a custom sink for log messages.
	// If nil, messages will be written to stdout.
	Logger *log.Logger
}

type Controller struct {
	sync.RWMutex
	conf         Configuration
	logger       *logger
	modbusClient *modbus.ModbusClient
	tags         []*Tag
	exit         bool

	// metrics
	errCounter *metrics.Counter
	reqCounter *metrics.Counter
}

func New(conf *Configuration) (c *Controller, err error) {
	defaults.SetDefaults(conf)
	c = &Controller{
		conf: *conf,
	}

	// Создаем метрики
	c.reqCounter = metrics.NewCounter("req_counter")
	c.errCounter = metrics.NewCounter("err_counter")

	// for an RTU over TCP device/bus (remote serial port or
	// simple TCP-to-serial bridge)
	c.modbusClient, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     c.conf.Url,
		Speed:   c.conf.Speed, // serial link speed
		Timeout: c.conf.Timeout,
	})
	if err != nil {
		return
	}

	err = c.modbusClient.SetUnitId(c.conf.DeviceId)
	if err != nil {
		return
	}

	err = c.modbusClient.Open()

	return
}

func (c *Controller) AddTag(tag *Tag) {
	c.Lock()
	defer c.Unlock()

	tag.Gauge = metrics.NewGauge(tag.Name, func() float64 {
		c.RLock()
		defer c.RUnlock()
		if tag.LastValue != nil {
			switch tag.Method {
			case READ_UINT:
				return float64(tag.LastValue.(uint16))
			case READ_FLOAT:
				return float64(tag.LastValue.(float32))
			}
		}
		return 0.0
	})

	c.tags = append(c.tags, tag)
}

func (c *Controller) Close() {
	c.exit = true
}

func (c *Controller) incCounter() {
	c.reqCounter.Inc()
}

func (c *Controller) incErrCounter() {
	c.errCounter.Inc()
}

func (c *Controller) Poll() {
	log.Println("Start polling...")

	c.exit = false
	needRestart := false
	for {
		if c.exit {
			break
		}
		for i, tag := range c.tags {
			// Принудительный рестарт
			if needRestart {
				log.Println("Restarting connect...")
				err := c.modbusClient.Open()
				if err != nil {
					log.Println("Can not open connect")
					break
				}
				needRestart = false
			}

			time.Sleep(c.conf.ReadPeriod)

			c.Lock()
			var err error
			var val interface{}

			switch tag.Method {
			case READ_UINT:
				val, err = c.modbusClient.ReadRegister(tag.Address, modbus.HOLDING_REGISTER)
				c.incCounter()
			case READ_FLOAT:
				val, err = c.modbusClient.ReadFloat32(tag.Address, modbus.HOLDING_REGISTER)
				c.incCounter()
			}

			// Обработка ошибок
			if err != nil {
				c.incErrCounter()
				if c.logger != nil {
					c.logger.customLogger.Print("Err get tag " + tag.Name + " err: " + err.Error())
				} else {
					log.Println("Err get tag " + tag.Name + " err: " + err.Error())
				}

				if cause, ok := err.(interface{ Unwrap() error }); ok {
					if _, ok := cause.(net.Error); ok {
						needRestart = true
						c.Unlock()
						break
					}
				}
				c.Unlock()
				continue
			}
			tag.Action(val, c.tags[i])
			c.Unlock()
		}
		time.Sleep(c.conf.PollingTime)
	}

	log.Println("End polling")
	err := c.modbusClient.Close()
	if err != nil {
		log.Println("Controller close error: " + err.Error())
	}
}
