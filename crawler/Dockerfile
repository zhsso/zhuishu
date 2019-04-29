from python:3-alpine
RUN apk add gcc
RUN pip install requests
RUN pip install beautifulsoup4
RUN pip install redis
RUN pip install bottle
RUN apk add musl-dev
RUN pip install gevent
COPY entrypoint.sh /entrypoint.sh
EXPOSE 8080
CMD ["/entrypoint.sh"]
