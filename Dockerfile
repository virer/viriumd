FROM golang:1.23.7

WORKDIR /

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /virium

EXPOSE 8787
ENV VG_NAME vg_data

# Run
CMD ["/virium"]