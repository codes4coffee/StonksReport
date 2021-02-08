FROM golang

WORKDIR /go/src/app
COPY . .

RUN go get -d -v
RUN go install -v
RUN apt-get install -y tzdata

CMD ["stonksReport"]