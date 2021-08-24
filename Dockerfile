FROM golang:1.17 AS builder

WORKDIR /avakian

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

RUN git describe >> tag

RUN echo "$(cat tag)"

RUN go build -ldflags="-X github.com/holedaemon/avakian/internal/version.version=$(cat tag)" -o avakian

FROM gcr.io/distroless/base:nonroot
COPY --from=builder /avakian/avakian /avakian
ENTRYPOINT [ "/avakian" ]