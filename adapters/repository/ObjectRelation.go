package repository

import (
	"fmt"
	"hoho-framework-v2/adapters/request"
	"hoho-framework-v2/infrastructure/kafka"
	"hoho-framework-v2/library/auth"
	objectrelation "hoho-framework-v2/usecase/object_relation"
	"os"
)

type objectRelation struct{}

func NewObjectRelationRepo() objectrelation.IObjectRelation {
	return &objectRelation{}
}

func (o *objectRelation) SaveRelations(authObject auth.AuthObject, nodes []objectrelation.Node, links []objectrelation.Link, host string) error {
	data := map[string]interface{}{
		"relations": links,
		"nodes":     nodes,
		"host":      host,
	}
	return kafka.Publish("object-relation-data", "save", data, authObject.GetUserTenantId())
}

func (o *objectRelation) DeleteNodesAndLinks(authObject auth.AuthObject, host string, deleteByCohost bool) error {
	url := fmt.Sprintf("%s/object-host/%s", os.Getenv("OBJECT_RELATION"), host)
	if deleteByCohost {
		url = fmt.Sprintf("%s/object-cohost/%s", os.Getenv("OBJECT_RELATION"), host)
	}
	res, err := request.Make(url).SetHeaders(map[string]string{"Authorization": authObject.GetToken()}).Delete()
	if err != nil {
		return err
	}
	if res.Status != 200 {
		return fmt.Errorf("DeleteNodesAndLinks: %s", string(res.Data.([]byte)))
	}
	return nil
}
