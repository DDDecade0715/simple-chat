package initialize

import (
	"gin-derived/global"
	es "gin-derived/pkg/elasticsearch"
)

func InitElasticsearch() *es.Elasticsearch {
	//实例一个Elasticsearch
	host := global.GCONFIG.Es.Host
	EsClient := es.GetEs("default", host)

	return EsClient
}
