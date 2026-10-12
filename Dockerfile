FROM docker.friedmicro-lab.org/base

RUN apt install nginx -y
COPY web/nginx.conf /etc/nginx/conf.d/default.conf

COPY web web
RUN cd web; npm ci; npm run build; mv dist/* /usr/share/nginx/html/
COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum
RUN go build cmd/api/main.go
RUN mv main planner

# Start both the Go API and nginx in the background.
CMD ["/bin/sh", "-c", "/planner & exec nginx -g \"daemon off;\""]

EXPOSE 5005
