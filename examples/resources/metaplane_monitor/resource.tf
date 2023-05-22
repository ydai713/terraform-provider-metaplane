data "metaplane_connection" "snowflake" {
  name = "snowflake"
}

resource "metaplane_monitor" "monitor" {
  absolute_path = "DATABASE.SCHEMA.TABLE"
  entity_type   = "TABLE"
  type          = "ROW_COUNT"
  cron_tab      = "* 2 * * *"
  connection_id = data.metaplane_connection.snowflake.id
}
