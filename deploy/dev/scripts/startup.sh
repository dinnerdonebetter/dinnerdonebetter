function start_watchtower() {
	docker run --detach \
	--name watchtower \
	--volume /var/run/docker.sock:/var/run/docker.sock \
	containrrr/watchtower
}

function start_caddy() {
	docker run --detach \
	--name caddy \
	--volume /home/sysadmin/Caddyfile:/etc/caddy/Caddyfile \
	--volume /home/sysadmin/caddy_data:/data \
	--publish 443:443 \
	caddy:2
}

function start_prixfixe() {
    docker run --detach \
		--name prixfixe \
		--volume /home/sysadmin/prixfixe.config.toml:/etc/prixfixe/config.toml \
		--env CONFIGURATION_FILEPATH=/etc/prixfixe/config.toml
		--publish 443:443 \
		registry.gitlab.com/prixfixe/prixfixe:dev
}
