FROM alpine:3.17
WORKDIR /app
ADD review review
EXPOSE 80
ENTRYPOINT [ "/app/review" ]