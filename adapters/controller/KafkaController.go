package controller

import (
	"errors"
	"hoho-framework-v2/infrastructure/kafka"
	"net/http"
)

type kafkaController struct {
}

type KafkaController interface {
	Subscribe(c *Context) error
	StopProcess(c *Context) error
}

func NewKafkaController() KafkaController {
	return &kafkaController{}
}

func (k *kafkaController) Subscribe(c *Context) error {
	// topics := c.FormValue("topic")
	// go k.consumer.Subscribe("","",)
	x := map[string]interface{}{"a": "s"}
	x1 := make([]map[string]interface{}, 0)
	x1 = append(x1, x)
	kafka.Publish("dssdsad", "log", x1, 0)
	return nil
}
func (k *kafkaController) StopProcess(c *Context) error {
	processId := c.Param("processId")
	if kafka.ConsumerRunning[processId] != nil {
		kafka.StopSubscribe(processId)
		return nil
	} else {
		c.Output(http.StatusNotFound, "", errors.New("process not running"))
		return errors.New("process not running")
	}
}
