FROM caddy:2-alpine

COPY Caddyfile /etc/caddy/Caddyfile

RUN mkdir -p /srv/.well-known

COPY apple-app-site-association /srv/.well-known/apple-app-site-association

EXPOSE 80
