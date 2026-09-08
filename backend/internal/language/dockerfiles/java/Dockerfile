FROM alpine:3.22 AS downloader
RUN apk add --no-cache wget && \
    wget -q -O /json.jar \
      https://repo1.maven.org/maven2/org/json/json/20260814/json-20260814.jar
FROM eclipse-temurin:21-jdk-alpine
COPY --from=downloader /json.jar /opt/lib/json.jar
WORKDIR /code
