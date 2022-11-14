FROM golang:1.19
WORKDIR /usr/src
COPY . ./
RUN go mod download
RUN go build -o ./app
RUN apt-get update
RUN apt install -y httpie
RUN apt-get install -y nginx
COPY default /etc/nginx/sites-available
RUN  apt-get update && apt-get install -y redis-server
CMD  service nginx start && service redis-server start && ./app && sleep infinity