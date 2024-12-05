schema "main" {
}

table "magnet_cache" {
  schema = schema.main

  column "store" {
    null = false
    type = varchar
  }
  column "hash" {
    null = false
    type = varchar
  }
  column "is_cached" {
    null = false
    type = bool
    default = false
  }
  column "modified_at" {
    null = false
    type = datetime
    default = sql("unixepoch()")
  }
  column "files" {
    null = false
    type = json
    default = sql("json('[]')")
  }
  primary_key {
    columns = [column.store, column.hash]
  }
}