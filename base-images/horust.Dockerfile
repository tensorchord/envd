FROM rust:bullseye as builder

RUN git clone https://github.com/FedericoPonzi/Horust.git /app --depth=1
WORKDIR /app
RUN cargo build --release

FROM scratch
COPY --from=builder /app/target/release/horust .
ENTRYPOINT [ "./horust" ]
