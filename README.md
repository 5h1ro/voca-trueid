# INSTALL

```
 1. docker pull ghcr.io/go-rod/rod:latest
 2. docker run --name browser -p 7317:7317 ghcr.io/go-rod/rod:latest
 3. docker pull redis
 4. docker run --name redis -p 6379:6379 -d redis
 5. docker build -t voca-trueid .
 6. docker run --name voca-trueid --network=host -p 9000:9000 voca-trueid
```
