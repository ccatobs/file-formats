FROM golang:1

RUN go get -v \
    github.com/xitongsys/parquet-go/writer \
    github.com/xitongsys/parquet-go-source/local

RUN apt-get update && apt-get install -y \
    python3-pip

# upgrade pip, otherwise fastparquet fails to install
RUN pip3 install -U pip

RUN pip3 install \
    fastparquet \
    pandas \
    pyarrow

