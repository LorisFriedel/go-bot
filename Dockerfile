FROM iron/base

WORKDIR /app

# Add the binary
ADD bin/go-bot /app/
ENTRYPOINT ["./go-bot"]
