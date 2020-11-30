FROM alpine:3.12 AS rie-image

RUN apk add --no-cache curl && \
  cd / && \
  curl -sSL -O https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/download/v1.0/aws-lambda-rie && \
  chmod +x /aws-lambda-rie


FROM pdf2png

COPY --from=rie-image /aws-lambda-rie /usr/local/bin/

ENTRYPOINT ["aws-lambda-rie", "/app/main"]
