#!/usr/bin/env sh

# Migrate django database
echo Migrating database...
./manage.py migrate

# Start Gunicorn processes
echo Starting Gunicorn webserver...
exec gunicorn automarkwatched.wsgi \
     --bind 0.0.0.0:8000 \
     --workers 3