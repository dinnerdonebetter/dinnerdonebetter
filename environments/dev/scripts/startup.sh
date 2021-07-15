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

function fuck_docker() {
  docker stop $(docker ps --all --quiet)
  docker rm $(docker ps --all --quiet)
}

function start() {
  docker-compose --file $HOME_DIR/docker-compose.yaml up
}

function init_search_index() {
  docker run --volume /home/sysadmin/search_indices:/output --user `id -u` registry.gitlab.com/prixfixe/prixfixe:index_init /index_initializer --output /output/$1_index.bleve --type=$1 --db_connection="postgresql://prixfixe_dev:vfhfFBwoCoDWTY86bVYa9znk1xcp19IO@database.prixfixe.dev:25060/dev_prixfixe?sslmode=require" --db_type=postgres --deadline=30s
  chown sysadmin:sysadmin -R /home/sysadmin/search_indices/$1_index.bleve
}

function init_all_indices() {
        init_search_index valid_ingredients
        init_search_index valid_instruments
        init_search_index valid_preparations
}
