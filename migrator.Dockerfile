# Stage 1 - downloading goose
FROM alpine:3.20 AS goose-downloader

ADD https://github.com/pressly/goose/releases/download/v3.23.0/goose_linux_x86_64 /bin/goose
RUN chmod +x /bin/goose

# Final stage
FROM alpine:3.20

# Copy goose from the first stage
COPY --from=goose-downloader /bin/goose /bin/goose

# Install bash and clean cache in one layer
RUN apk add --no-cache bash

WORKDIR /opt/app

# Copy migration files and script
COPY migrations/*.sql migrations/
COPY migration.sh ./migration.sh
RUN chmod +x migration.sh

ENTRYPOINT ["bash", "migration.sh"]
