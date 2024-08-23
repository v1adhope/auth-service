from busybox:1.36.1
label org.opencontainers.image.authors="Vladislav Gardner <vladislavgardner@gmail.com>"

workdir service

copy .bin/auth-service .env ./

user 1000:1000

cmd ["./auth-service"]
