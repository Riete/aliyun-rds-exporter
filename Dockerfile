FROM riet/golang:1.13.10 as backend
COPY . .
RUN unset GOPATH && go build -mod=vendor

FROM riet/centos:7.4.1708-cnzone
COPY --from=backend /go/aliyun-rds-exporter /opt/aliyun-rds-exporter
EXPOSE 10001
CMD /opt/aliyun-rds-exporter
