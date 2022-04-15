# start app
.PHONY : run
run:
	go mod download && go run main.go