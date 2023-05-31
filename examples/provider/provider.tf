terraform {
  required_providers {
    metaplane = {
      source  = "ydai713/metaplane"
      version = "0.0.5"
    }
  }
}

provider "metaplane" {
  api_key = ""
}
