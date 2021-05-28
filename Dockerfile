# Set the base image
FROM golang:1.16 

LABEL maintainer="ssong. <mwsong@rockplace.co.kr>"

# Add ghost user
RUN useradd -m -s /bin/bash ghost && \
    mkdir -p /usr/app/jpetstore && \
    chown -R ghost:ghost /usr/app/jpetstore 

COPY ./ /usr/app/jpetstore
 
# Switch to ghost user
USER ghost:ghost
WORKDIR /usr/app/jpetstore 

RUN go build -o run main.go && echo "go run main.go" > /usr/app/jpetstore/entrypoint.sh && \ 
    chmod 700 /usr/app/jpetstore/entrypoint.sh 

#ENTRYPOINT ["bash","-c","./entrypoint.sh"]
ENTRYPOINT ["bash","-c","/usr/app/jpetstore/run"]