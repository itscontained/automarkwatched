# AutoMarkWatched
[![Documentation Status](https://readthedocs.org/projects/plexshowsilencer/badge/?version=latest)](http://plexshowsilencer.readthedocs.io/en/latest/?badge=latest)

A web application for marking TV shows as watched/unwatched for [Plex Media Server](https://plex.tv).

AutoMarkWatched runs on [Django](https://www.djangoproject.com/) and uses sqlite3 as a database.

Documentation hasnt started yet but will be hosted at [Read the Docs](http://plexshowsilencer.readthedocs.io/en/latest/)

# Quickstart
(Until i have time to document in more detail)
1. Git clone repo
2. Generate secret_key `python3 generate_secret_key.py`
3. Edit automarkwatched/settings.py
   * Add generated key to SECRET_KEY = '<here>'
   * Add ip addres or fqdn to ALLOWED_HOSTS
4. Create database `python3 manage.py migrate`
5. Run server `python3 manage.py runserver 0:8000`

Note: If you know how to use gunicorn and nginx, i have configs for both in the root dir. Will add to docs eventually 

![Home Page](https://imgur.com/a/JeJQj "Home Page")
![Bulk Edit](https://imgur.com/a/8qZgq "Bulk Edit")
![Settings](https://imgur.com/a/6OzT6 "Settings")
