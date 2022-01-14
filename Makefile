EXEC = hwbot

all:
	go build -o ./bin/$(EXEC) ./

brun:
	go build -o ./bin/$(EXEC) ./ && ./bin/$(EXEC)

run:
	./bin/$(EXEC)
