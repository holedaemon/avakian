FROM golang:1.16 AS builder

WORKDIR /avakian

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN go build -o avakian

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /avakian/avakian /avakian
ENTRYPOINT [ "/avakian" ]