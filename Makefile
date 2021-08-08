PLUGINS = plugins
EXEC = hwbot

all: plugins
	go build -o ./bin/$(EXEC) ./

plugins:
	go build -buildmode=plugin -o bin/modules/debug.bmod modules/debug/debug.go
	go build -buildmode=plugin -o bin/modules/yalm.bmod modules/yalm/yalm.go
# go build -o ./bin/$(PLUGINS)/ -buildmode=plugin ./plugins/*/*.go

brun:
	go build -o ./bin/$(EXEC) ./ && ./bin/$(EXEC)

run:
	./bin/$(EXEC)
