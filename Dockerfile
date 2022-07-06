FROM golang:1.17-alpine as base

FROM base as frontendbuild
WORKDIR /app
RUN apk add --update nodejs npm
COPY ./frontend/package.json /app
COPY ./frontend/package-lock.json /app
RUN npm install

COPY ./frontend /app
RUN npm run build

FROM base as gobuild
WORKDIR /app
COPY go.sum .
COPY go.mod .
# RUN go mod download
COPY . .
COPY --from=frontendbuild /app/dist /app/server
RUN apk add --update build-base
RUN CGO_ENABLED=1 go build -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=gobuild /app/bookcatalog /app/bookcatalog
COPY --from=gobuild /app/db/migrations /app/db

# TODO: dont use hardcoded migrations path
# TODO: Maybe just use the entrypoint command and parameterize user input
CMD ["/app/bookcatalog", "start", "-d", "/app/db/migrations", "-l", "/mount/library", "-i", "/mount/images" ]

