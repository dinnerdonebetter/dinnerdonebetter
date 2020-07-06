# deployment

## operating systems

The only acceptable operating systems to use on servers are:

1. Debian 10
1. Ubuntu 20.04

## the dev environment in a nutshell

There are three containers running on a DigitalOcean VPS, which serves at https://prixfixe.dev 

1. [Caddy](https://caddyserver.com/) - a reverse proxy in front of the PrixFixe server. Caddy also ensures we serve over HTTPS, a requirement of the domain TLD.
1. [PrixFixe](https://gitlab.com/prixfixe/prixfixe) - the service we all know and love.
1. [Watchtower](https://github.com/containrrr/watchtower) - which checks to see if there are any new versions of the Caddy and PrixFixe containers, and gracefully pulls/restarts the relevant container.

## how to replicate the dev environment

1. Create infrastructure in DigitalOcean
    1. A managed PostgresSQL database (DO calls them clusters regardless of size)
    1. A VPS (droplet) with appropriate resources. Currently `dev` is equipped with: 
        - Ubuntu 20.04
        - 1 vCPU
        - 2 GB RAM
        - 50 GB ssd disk
1. Create the DNS records in Cloudflare
    1. Create a `CNAME` record to point to the domain of the new Postgres cluster
    1. Create an `A` record in Cloudflare to point to the new droplet's IPv4 address
    1. Create an `AAAA` record in Cloudflare to point to the new droplet's IPv6 address
1. VPS setup
    1. Create a user called `sysadmin` and save the password in 1Password
    1. `usermod -aG sudo sysadmin`
    1. Ensure the relevant SSH public keys are present in `/home/sysadmin/.ssh/authorized_keys` and that the file belongs to the `sysadmin` user and group
    1. Exit root session
    1. Log in as the `sysadmin` user 
    1. Install `docker` and `docker-compose`
    1. Run `docker login registry.gitlab.com --username $GITLAB_USERNAME --password $GITHUB_ACCESS_TOKEN`
    1. Copy these files to the appropriate locations:
        1. `$(pwd)/docker-compose.yaml` ➡️ `/home/sysadmin/docker-compose.yaml`
        1. `$(pwd)/config.toml` ➡️ `/home/sysadmin/prixfixe.config.toml`
        1. `$(pwd)/caddy/Caddyfile` ➡️ `/home/sysadmin/Caddyfile`
        1. `$(pwd)/scripts` ➡️ `/home/sysadmin/.scripts.sh`
    1. Edit `/home/sysadmin/.bashrc` so that it runs `. /home/sysadmin/.scripts.sh`
    1. Run `source /home/sysadmin/.scripts.sh` as the logged in sysadmin user
    1. Run `start`
    
