FROM ubuntu
ARG TOKEN
ENV BOT_TOKEN=${TOKEN}
#INSTALL CA-CERTS
RUN apt -q update && apt install -qy ca-certificates
RUN update-ca-certificates
#COPY BINARY
COPY main .
CMD ["./main"]