FROM ubuntu
RUN apt update
RUN apt install -y ca-certificates
COPY entrypoint.sh /entrypoint.sh
CMD ["/entrypoint.sh"]
