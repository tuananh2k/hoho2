package objectrelation

import "hoho-framework-v2/library/auth"

type Node struct {
	Name  string `json:"name"`
	Id    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
	Host  string `json:"host"`
}

type Link struct {
	Start string `json:"start"`
	End   string `json:"end"`
	Type  string `json:"type"`
	Host  string `json:"host"`
}

type IObjectRelation interface {
	SaveRelations(authObject auth.AuthObject, node []Node, links []Link, host string) error
	DeleteNodesAndLinks(authObject auth.AuthObject, host string, deleteByCohost bool) error
}
