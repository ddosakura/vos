dist/vos: vos/*.go
	go build -o dist/vos ./vos

clean:
	rm ./**/*.log & rm ./**/*_history
