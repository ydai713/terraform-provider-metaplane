---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "metaplane Provider"
subcategory: ""
description: |-
  
---

# metaplane Provider



## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `api_key` (String) Metaplane API Key
