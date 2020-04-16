#!/usr/bin/env sh

# Migrate django database
echo Migrating database...
./manage.py migrate

# Collect static files
echo Collecting static files...
./manage.py collectstatic

# Start Gunicorn processes
echo Starting webserver...
exec runserver 0:8000
