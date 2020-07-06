# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = var.do_token
}

data "digitalocean_database_cluster" "database" {
  name = "database"
}

data "digitalocean_droplet" "dev_server" {
  name = "prixfixe-dev"
}
