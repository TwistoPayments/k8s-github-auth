FROM alpine:3.9

RUN apk --no-cache --update add ca-certificates
COPY build/main /k8s-git-auth
CMD ["/k8s-git-auth"]
