
go-analytics
============

A analytics system write by golang.

## Load template to ElasticSearch
~~~bash
curl -XPOST http://localhost:9200/analytics/ -d@index.template.json
~~~
make sure return `{"acknowledged":true}`

## Run
~~~bash
export ENV_ELASTICSEARCH_HOST=localhost
go run server.go :8001
~~~

## Test
~~~bash
export ENV_TEST_SERVER=http://localhost:8001
go test
~~~
