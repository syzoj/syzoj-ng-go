FROM golang:1.13
WORKDIR /go/src/github.com/syzoj/syzoj-ng-go/
COPY . .
RUN go get -v github.com/syzoj/syzoj-ng-go/svc/main/app && go get -v github.com/syzoj/syzoj-ng-go/svc/main/migrate

FROM ubuntu
WORKDIR /app
COPY --from=0 /go/bin/ .
