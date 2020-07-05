# basics
HOME_DIR=/home/sysadmin

# reverse proxy things
REVERSE_PROXY_CONTAINER_NAME=caddy
REVERSE_PROXY_CONFIG_PATH=$HOME_DIR/Caddyfile

# misc ops things
CONTAINER_NETWORK_NAME=dev_network
WATCHTOWER_IMAGE_NAME=watchtower

# service things
APP_CONTAINER_NAME=prixfixe
APP_CONTAINER_IMAGE=registry.gitlab.com/prixfixe/prixfixe:dev
APP_CONFIG_PATH=$HOME_DIR/prixfixe.config.toml

# container startup functions
function start_network(){
  if [[ $(docker network ls --filter name=$CONTAINER_NETWORK_NAME | grep "" -c) -eq "2" ]]; then
    echo "docker network already exists"
  else
    docker network create --attachable $CONTAINER_NETWORK_NAME
  fi
}

function start_watchtower() {
	docker run \
	--detach \
	--name $WATCHTOWER_IMAGE_NAME \
	--volume $HOME_DIR/.docker/config.json:/config.json \
	--volume /var/run/docker.sock:/var/run/docker.sock \
	containrrr/watchtower:latest
}

function start_caddy() {
	docker run \
	--detach \
	--name $REVERSE_PROXY_CONTAINER_NAME \
	--volume REVERSE_PROXY_CONFIG_PATH:/etc/caddy/Caddyfile \
	--volume $HOME_DIR/caddy_data:/data \
	--publish 443:443 \
	--network $CONTAINER_NETWORK_NAME \
	--network-alias $REVERSE_PROXY_CONTAINER_NAME \
	caddy:2
}

function start_prixfixe() {
	docker run --detach \
	--restart always \
	--name $APP_CONTAINER_NAME \
	--env CONFIGURATION_FILEPATH=/etc/prixfixe/config.toml \
	--volume $APP_CONFIG_PATH:/etc/prixfixe/config.toml \
	--publish 8888:8888 \
	--network $CONTAINER_NETWORK_NAME \
	--network-alias prixfixe_server \
	$APP_CONTAINER_IMAGE
}

function edit_scripts() {
	vim $HOME_DIR/.scripts.sh && source $HOME_DIR/.bashrc
}

function edit_caddyfile() {
	vim $REVERSE_PROXY_CONFIG_PATH && docker restart $REVERSE_PROXY_CONTAINER_NAME && sleep 2 && docker ps
}

function edit_server_config() {
	vim $APP_CONFIG_PATH && docker restart $APP_CONTAINER_NAME && sleep 2 && docker ps
}

function wipe_docker() {
	docker rm -f $REVERSE_PROXY_CONTAINER_NAME $APP_CONTAINER_NAME $CONTAINER_NETWORK_NAME
}

function spin_up() {
  start_network &&
  start_watchtower &&
  start_caddy &&
  start_prixfixe
}