FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends g++ wget ca-certificates && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /opt/json
RUN wget -O json.hpp https://github.com/nlohmann/json/releases/latest/download/json.hpp && \
    g++ -std=c++17 json.hpp -o json.hpp.gch
WORKDIR /code