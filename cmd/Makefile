test:
	go test -v .././...

compile:
	echo "Compiling for Linux..."
	GOOS=linux go build -o ../bin/idealista-api;

run-api: compile
	cd .. && ./bin/idealista-api

run-reports: compile
	cd .. && ./bin/idealista-api sendMonthlyReports

run-add-flats: compile
	cd .. && ./bin/idealista-api addFlats

all: test compile