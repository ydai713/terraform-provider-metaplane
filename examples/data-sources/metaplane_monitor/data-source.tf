data "metaplane_connection" "snowflake" {
  name = "snowflake"
}

data "metaplane_monitor" "monitor" {
  connection_id = data.metaplane_connection.snowflake.id
  monitor_id    = "string"
}
