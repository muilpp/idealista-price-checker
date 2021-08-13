test:
	go test -v .././...

compile:
	echo "Compiling for Linux..."
	GOOS=linux go build -o ../bin/idealista-api;

run-api:
	cd .. && ./bin/idealista-api

run-reports:
	cd .. && ./bin/idealista-api sendMonthlyReports

all: test compile