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
# RUN CGO_ENABLED=1 go build -o bookcatalog . # -ldflags='-extldflags=-static'
RUN CGO_ENABLED=1 go build -a -ldflags '-linkmode external -extldflags "-static"' .

FROM scratch
COPY --from=gobuild /app/bookcatalog /app/bookcatalog

WORKDIR /mount
CMD ["/app/bookcatalog", "start" ]

