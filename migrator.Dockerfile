# Set the base image to Alpine Linux 3.20
FROM alpine:3.20

# Update and upgrade the package index, install Bash, and remove cached packages
RUN apk update && \
    apk upgrade && \
    apk add bash && \
    rm -rf /var/cache/apk/*

# Download the Goose binary and add it to the /bin/ directory
ADD https://github.com/pressly/goose/releases/download/v3.23.0/goose_linux_x86_64 /bin/goose
# Make the Goose binary executable
RUN chmod +x /bin/goose

# Set the working directory to /opt/app
WORKDIR /opt/app

# Copy the SQL migration files, the migration script, and the environment file to the container
COPY migrations/*.sql migrations/
COPY migration.sh ./migration.sh

# Make the migration script executable
RUN chmod +x migration.sh

# Set the default command to run the migration script
ENTRYPOINT ["bash", "migration.sh"]
