
.PHONY: run
run:
	docker build -t parquet-file-format-test .
	docker run --rm -it -v "$$PWD:/stuff" -w /stuff parquet-file-format-test ./run.sh

