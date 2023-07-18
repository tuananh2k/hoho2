FROM golang:1.19.10
RUN mkdir /app
RUN mkdir /src
WORKDIR /app
ADD ./ ./

RUN go mod download
RUN go build -o main
RUN chmod -R 777 /app/log
RUN mkdir /app/request-log
RUN chmod -R 777 /app/request-log
EXPOSE 1323
WORKDIR /src
CMD [ "/app/main" ]