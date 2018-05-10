all:	clcheat clcheat.exe

clcheat:
	GOOS=linux go build

clcheat.exe:
	GOOS=windows go build

clean:
	rm -f clcheat clcheat.exe clcheat.log
