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
	--network dev_network \
	--network-alias caddy
	caddy:2
}

function start_network(){
	docker network create --attachable dev_network
}

function reset_docker() {
	docker rm -f caddy prixfixe
}

function edit_scripts() {
	vim /home/sysadmin/.scripts.sh && source /home/sysadmin/.bashrc
}

function edit_caddyfile() {
	vim /home/sysadmin/Caddyfile && docker restart caddy && sleep 1 && docker ps
}

function start_prixfixe() {
	docker run --detach \
		--name prixfixe \
		--volume /home/sysadmin/prixfixe.config.toml:/etc/prixfixe/config.toml \
		--env CONFIGURATION_FILEPATH=/etc/prixfixe/config.toml \
		--publish 8888:8888 \
		--restart=always \
		--network dev_network \
		--network-alias prixfixe_server
		registry.gitlab.com/prixfixe/prixfixe:dev
}
