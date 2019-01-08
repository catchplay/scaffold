
#build stage
FROM golang:alpine AS builder
WORKDIR /go/src/{{.ProjectPath}}
COPY . .
RUN apk add --no-cache git
RUN go get ./...
RUN go build -o {{.ProjectName}} ./cmd/main.go

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk add --no-cache bash
COPY --from=builder /go/src/{{.ProjectPath}}/{{.ProjectName}} /{{.ProjectName}}
COPY --from=builder /go/src/{{.ProjectPath}}/config/config.yml /config.yml

ADD https://raw.githubusercontent.com/vishnubob/wait-for-it/master/wait-for-it.sh ./wait-for-it.sh
RUN ["chmod", "+x", "./wait-for-it.sh"]
LABEL Name={{.ProjectName}}
EXPOSE 8080
