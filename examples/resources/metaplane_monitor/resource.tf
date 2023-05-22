terraform {
  required_providers {
    metaplane = {
      source = "registry.terraform.io/klaviyo/metaplane"
    }
  }
}

provider "metaplane" {
}

data "metaplane_connection" "snowflake" {
  name = "Klaviyo Prod"
}

resource "metaplane_monitor" "monitor" {
  absolute_path = "AIRBYTE.STRIPE.CUSTOMERS"
  entity_type   = "TABLE"
  type          = "ROW_COUNT"
  cron_tab      = "* 2 * * *"
  connection_id = data.metaplane_connection.snowflake.id
}

resource "metaplane_monitor" "monitor_2" {
  absolute_path = "AIRBYTE.STRIPE.SUBSCRIPTIONS"
  entity_type   = "TABLE"
  type          = "ROW_COUNT"
  cron_tab      = "* 1 * * *"
  connection_id = data.metaplane_connection.snowflake.id
}

resource "metaplane_monitor" "monitor_3" {
  absolute_path = "KLAVIYO.SALESFORCE.ACCOUNT"
  entity_type   = "TABLE"
  type          = "ROW_COUNT"
  cron_tab      = "0 * * * *"
  connection_id = data.metaplane_connection.snowflake.id
}

resource "metaplane_monitor" "monitor_4" {
  absolute_path = "KLAVIYO.SALESFORCE.USER"
  entity_type   = "ROW_COUNT"
  type          = "FRESHNESS"
  cron_tab      = "0 * * * *"
  connection_id = data.metaplane_connection.snowflake.id
}
