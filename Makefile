GOROOT=${HOME}/Downloads/go

all:
	~/gopath/bin/statik
	gofmt -w=true lws.go
	#go build lws.go
	#./lws
	GOOS=linux GOARCH=arm ${GOROOT}/bin/go \
		build -o lws-arm lws.go
	adb push lws-arm /data/lws
	adb shell /data/lws


#https://github.com/syncthing/syncthing.git

deps:
	go get github.com/goji/httpauth
	go get github.com/rakyll/statik

clean:
	rm -f lws lws-arm
