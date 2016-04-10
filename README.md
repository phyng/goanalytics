
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
go run server.go 8001
~~~

## Test
~~~bash
go run test.go 8001
~~~
