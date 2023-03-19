FROM alpine:3.17
ADD review /app/review
EXPOSE 80
ENTRYPOINT [ "/app/review" ]