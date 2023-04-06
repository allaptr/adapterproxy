FROM golang:1.19-alpine

WORKDIR /go/work

COPY go.* ./
RUN go mod download

COPY cache/*.go ./cache/
COPY model/*.go ./model/
COPY cmd/backendify/*.go ./

RUN go build -o ./backendify

EXPOSE 9000

ENTRYPOINT [ "./backendify", "-status", "200"]