FROM eclipse-temurin:21-jdk-alpine
RUN apk add --no-cache wget
WORKDIR /opt/lib
RUN wget -O json.jar https://repo1.maven.org/maven2/org/json/json/20260814/json-20260814.jar && \
    apk del wget
WORKDIR /code
