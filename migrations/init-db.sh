#!/bin/bash
set -e

echo 'Migrating ...'
psql --username $POSTGRES_USER -d $POSTGRES_DB < /docker-entrypoint-initdb.d/20170518032523_create_user_table.up.sql
