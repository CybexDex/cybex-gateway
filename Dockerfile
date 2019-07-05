FROM golang:1.12

RUN apt-get update && apt-get install -y net-tools vim telnet

LABEL gateway.version=$VERSION
LABEL gateway.build_date=$BUILD_DATE

WORKDIR /usr/src/app

ADD . .
EXPOSE 8081 8182

# make sure gateway server start after postgres service
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.5.0/wait /wait
RUN chmod +x /wait

CMD ["sh", "-c", "/wait && GO111MODULE=on env=dev ./entry.sh"]
