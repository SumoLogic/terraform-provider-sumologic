terraform {
  required_providers {
    sumologic = {
      #source = "sumologic/sumologic"
      #version = "~> 2.19.1"
      source = "sumologic.com/dev/sumologic"
      version = "~> 1.0.0"
    }
  }
}

provider "sumologic" {
    base_url    = "https://stag-api.sumologic.net/api/"
    access_id   = "sulztB2SUcqONS"
    access_key  = "DXrQdslKmGYX5u1XnAKe9r23C8xaK9a47ck4vuxQFNmrjtazkEImX6V0SvXrsH6S"
}

