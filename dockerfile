FROM alpine:latest AS alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
WORKDIR /app
COPY ./build /app/build
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
CMD ["bash"]
ENTRYPOINT ["/app/build"]
