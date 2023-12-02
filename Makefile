compile = env GOOS=linux  GOARCH=arm64  go build -v -ldflags '-s -w -v' -tags lambda.norpc -o

build: gomodgen
	go mod download github.com/aws/aws-lambda-go
	go mod download github.com/sirupsen/logrus
	go get github.com/sirupsen/logrus@v1.9.3

	$(compile) bin/manageLogs/bootstrap manageLogs/levelLog.go
	$(compile) bin/cloudwatchlogs/bootstrap cloudwatchlogs/cloudwatchlogs.go

zip:
	zip -j -r bin/manageLogs.zip bin/manageLogs/bootstrap
	zip -j -r bin/cloudwatchlogs.zip bin/cloudwatchlogs/bootstrap
clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build zip
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
