data "metaplane_connection" "snowflake" {
  name = "snowflake"
}

resource "metaplane_monitor" "monitor" {
  absolute_path = "DATABASE.SCHEMA.TABLE"
  entity_type   = "TABLE"
  type          = "ROW_COUNT"
  cron_tab      = "* 2 * * *"
  connection_id = data.metaplane_connection.snowflake.id

  custom_sql              = ""
  custom_where_clause     = ""
  incremental_column_name = ""
  incremental_days        = 1
  incremental_hours       = 0
  incremental_minutes     = 0
}
