resource "digitalocean_app" "prixfixe" {
  spec {
    name   = "dev"
    region = "nyc3"

    domain {
      name = "api.prixfixe.dev"
    }

    service {
      name               = "api"
      environment_slug   = "go"
      instance_count     = 1
      instance_size_slug = "professional-xs"

      image {
        registry_type = "DOCR"
        registry      = ""
        repository    = "api_server"
      }

      http_port = 8000

      routes {
        path = "/"
      }

      run_command = "/server"
    }

    database {
      engine     = "PG"
      production = false
      db_name    = "dev-prixfixe"
      db_user    = "admin"
      version    = "12"
    }
  }

  depends_on = [
    digitalocean_container_registry.dev,
  ]
}