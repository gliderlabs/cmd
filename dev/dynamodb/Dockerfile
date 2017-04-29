FROM anapsix/alpine-java:8_server-jre
RUN \
    apk update \
    && apk add ca-certificates \
    && update-ca-certificates \
    && apk add openssl \
    && wget -q -O - http://dynamodb-local.s3-website-us-west-2.amazonaws.com/dynamodb_local_latest.tar.gz | tar xz
ENTRYPOINT ["/opt/jdk/bin/java", "-Djava.library.path=./DynamoDBLocal_lib", "-jar", "DynamoDBLocal.jar"]
CMD ["-help"]
