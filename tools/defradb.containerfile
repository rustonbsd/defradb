# syntax=docker/dockerfile:1

# An image to run defradb.

# Stage: build
# Several steps are involved to enable caching and because of the behavior of COPY regarding directories.
FROM docker.io/golang:1.24 AS build
WORKDIR /repo/
COPY go.mod go.sum Makefile ./
RUN make deps:modules
COPY . .
RUN make deps:playground
ENV BUILD_TAGS=playground,telemetry
ENV CGO_LDFLAGS="-L/lib -liroh"

# manually copy libwasmer.so to fix linking issue https://github.com/wasmerio/wasmer-go/issues/281
RUN export WASMER_ARCH=$(go env GOHOSTARCH | sed "s/arm64/aarch64/") && \
    export WASMER_PATH=$(go env GOMODCACHE)/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-$WASMER_ARCH/libwasmer.so && \
    cp $WASMER_PATH /lib/libwasmer.so
# manually copy libiroh.so to fix linking issue (dir is build with abs path, so put it in same place everywhere)
RUN export IROH_ARCH=$(go env GOHOSTARCH | sed "s/amd64/x86_64/" | sed "s/arm64/aarch64/") && \
    export IROH_PATH=$(go env GOMODCACHE)/github.com/rustonbsd/go-libp2p-iroh-transport@v0.1.1/ffi/libs/libiroh_${IROH_ARCH}_unknown_linux_gnu.so && \
    cp $IROH_PATH /lib/libiroh.so

RUN make build

# Stage: RUN
FROM debian:bookworm-slim
COPY --from=build /repo/build/defradb /defradb
COPY --from=build /lib/libwasmer.so /lib/libwasmer.so
COPY --from=build /lib/libiroh.so /lib/libiroh.so

# Documents which ports are normally used.
# To publish the ports: `docker run -p 9181:9181` ...
EXPOSE 9161
EXPOSE 9171
EXPOSE 9181

# Default command provided for convenience.
# e.g. docker run -p 9181:9181 source/defradb  start --url 0.0.0.0:9181
ENTRYPOINT [ "/defradb" ]
