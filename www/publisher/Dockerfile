FROM alpine:3.5
RUN apk --update add openssh
RUN ssh-keygen -t rsa -b 4096 -N "" -f /root/.ssh/id_rsa

RUN apk add bash hugo make git ca-certificates openssl
RUN update-ca-certificates
RUN git config --global user.name "Publisher"
RUN git config --global user.email "progrium+cmd@gmail.com"
RUN echo -e "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config

WORKDIR /work
COPY ./publish /publish
ENTRYPOINT ["/publish"]

RUN cat /root/.ssh/id_rsa.pub
