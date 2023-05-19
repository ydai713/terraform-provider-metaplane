data "metaplane_connection" "snowflake" {
  name = "Klaviyo Prod"
}

data "metaplane_monitor" "monitor" {
  connection_id = data.metaplane_connection.snowflake.id
  monitor_id    = "d8af199c-ab9f-4080-ac00-fab948fea0c9"
}
