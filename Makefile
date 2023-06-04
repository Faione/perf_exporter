APP=perf_exporter

BUILD_ENV=GOOS=linux CGO_ENABLED=0

build:  
	${BUILD_ENV} go build \
	   -trimpath \
	   -ldflags '-s -w' \
	   -o bin/${APP} .

clean:
	rm bin/*

image: build
	podman build -t perf_exporter:`./bin/${APP} -v | awk '{print $$3}'` .
   