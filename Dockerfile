FROM alpine:3.17
ADD review /app/review
EXPOSE 8000
ENTRYPOINT [ "/app/review" ]