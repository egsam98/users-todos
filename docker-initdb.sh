echo "create databases $POSTGRES_DEV, $POSTGRES_TEST"
psql -U postgres -c "create database $POSTGRES_DEV;"
psql -U postgres -c "create database $POSTGRES_TEST;"

SCHEMA_FOLDER=/docker-entrypoint-initdb.d/schema

echo "running scripts from $SCHEMA_FOLDER"
for filename in $SCHEMA_FOLDER/*.sql; do
  sql=$(cat $filename)
  sql=${sql%%;*}
  echo $sql
  psql -U postgres -d $POSTGRES_DEV -c "$sql"
done