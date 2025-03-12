
clean:
	rm -f bin/go_by_example
	rmdir bin

build:
	go build -o bin/go_by_example go_by_example.go

run1:
	@bin/go_by_example -word=opt -numb=7 -fork -svar=flag

run2:
	@bin/go_by_example -word=opt

run3:
	@bin/go_by_example -word=opt a1 a2 a3

run4:
	@bin/go_by_example -h

run5:
	@bin/go_by_example -wat
