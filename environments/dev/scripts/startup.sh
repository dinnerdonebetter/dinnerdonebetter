# basics
HOME_DIR=/home/sysadmin

# reverse proxy things
REVERSE_PROXY_CONTAINER_NAME=caddy
REVERSE_PROXY_CONFIG_PATH=$HOME_DIR/Caddyfile

# service things
APP_CONTAINER_NAME=prixfixe
APP_CONFIG_PATH=$HOME_DIR/prixfixe.config.toml

function edit_scripts() {
	vim $HOME_DIR/.scripts.sh && source $HOME_DIR/.bashrc
}

function edit_caddyfile() {
	vim $REVERSE_PROXY_CONFIG_PATH && docker restart $REVERSE_PROXY_CONTAINER_NAME && sleep 2 && docker ps
}

function edit_server_config() {
	vim $APP_CONFIG_PATH && docker restart $APP_CONTAINER_NAME && sleep 2 && docker ps
}

function start() {
  docker-compose --file $HOME_DIR/docker-compose.yaml up
}