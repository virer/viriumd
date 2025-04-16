
clean:
	rm -f bin/viriumd

mod-check:
	go mod verify && [ "$(shell sha512sum go.mod)" = "`sha512sum go.mod`" ] || ( echo "ERROR: go.mod was modified by 'go mod verify'" && false )


all:
	cd cmd; CGO_ENABLED=0 GOOS=linux go build -o ../bin/viriumd 
