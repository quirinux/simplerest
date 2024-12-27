FROM golang:latest AS builder
WORKDIR /src
COPY . /src
RUN make -C /src

FROM scratch AS runner
COPY --from=builder /src/bin/simplerest /bin
WORKDIR /app
EXPOSE 8888/tcp
CMD ["/bin/simplerest", "-config", "/app/simplerest.toml"]
