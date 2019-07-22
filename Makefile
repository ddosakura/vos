.PHONY: build
build: dist/vterm/pure dist/vos-base # dist/vos-fs dist/vos-script
	echo "Success"

dist/vterm/pure: vterm/pure/*.go proto/auth/auth.pb.go
	go build -o dist/vterm ./vterm/pure

dist/vos-base: *.go example/base/*.go proto/auth/auth.pb.go
	go build -o dist/vos-base ./example/base
dist/vos-fs: *.go example/fs/*.go proto/auth/auth.pb.go
	go build -o dist/vos-fs ./example/fs
dist/vos-script: *.go example/script/*.go proto/auth/auth.pb.go
	go build -o dist/vos-script ./example/script

proto/auth/auth.pb.go: proto/auth/auth.proto
	protoc --go_out=. proto/auth/auth.proto

.PHONY: clean
clean:
	rm ./**/*.log & rm ./**/*_history
	rm proto/*/*.pb.go
	rm ./**/*.sock
