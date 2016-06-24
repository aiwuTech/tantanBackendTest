# tantanBackendTest
Tantan Back-End Developer Test

## Install
    go get -u -v github.com/aiwuTech/tantanBackendTest
    cd $GOPATH/src/github.com/aiwuTech/tantanBackendTest
    go build

## run
1.you can simple run using the command(make sure ):
    ./tantanBackendTest -c ./config.json

2.or using command: ./tantanBackendTest --help for more details

## configuration
see config.example.json for more config, using -c args specify the location of config.json


## NOTICE
1. require postgreSQL for db, but i am using [gorm](https://github.com/jinzhu/gorm) for database support, and in this case mysql is used. But gorm also support postgreSQL, so i think no big deal.
2. at the pdf, [mux](https://github.com/gorilla/mux) is recommend, but i am familiar with [martini](https://github.com/go-martini/martini)