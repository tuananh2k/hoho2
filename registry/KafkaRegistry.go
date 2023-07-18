package registry

import (
	"hoho-framework-v2/adapters/controller"
)

func (r *registry) NewKafkaController() controller.KafkaController {
	return controller.NewKafkaController()
}
