FROM docker.friedmicro-lab.org/base AS build

COPY web web
RUN cd web; npm ci; npm run build; mv dist /opt/
COPY cmd cmd
COPY internal internal
COPY go.mod go.mod
COPY go.sum go.sum
RUN go build cmd/api/main.go
RUN mv main /opt/planner

FROM docker.friedmicro-lab.org/base-nginx
COPY web/nginx.conf /etc/nginx/
COPY --from=build /opt/planner /opt/planner
COPY --from=build /opt/dist/ /htmlcontent/

# Start both the Go API and nginx in the background.
CMD ["/bin/sh", "-c", "/opt/planner & exec nginx -g \"daemon off;\""]

EXPOSE 5005
