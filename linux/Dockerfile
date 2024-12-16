FROM ubuntu:22.04

# Set non-interactive frontend to suppress prompts
ENV DEBIAN_FRONTEND=noninteractive

# Install essential tools and systemd
RUN apt-get update && apt-get install -y \
    systemd \
    systemd-sysv \
    openssh-server \
    curl \
    vim \
    git \
    && apt-get clean && rm -rf /var/lib/apt/lists/*

# Set up SSH
RUN mkdir /var/run/sshd && echo 'root:password' | chpasswd

# Expose ports for SSH
EXPOSE 22

# Start systemd
# CMD ["/lib/systemd/systemd"]
CMD ["tail", "-f", "/dev/null"]