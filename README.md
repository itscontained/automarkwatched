# AutoMarkWatched
[![Documentation Status](https://readthedocs.org/projects/plexshowsilencer/badge/?version=latest)](http://plexshowsilencer.readthedocs.io/en/latest/?badge=latest)

A web application for marking TV shows as watched/unwatched for [Plex Media Server](https://plex.tv).

AutoMarkWatched runs on [Django](https://www.djangoproject.com/) and uses sqlite3 as a database.

Documentation hasnt started yet but will be hosted at [Read the Docs](http://plexshowsilencer.readthedocs.io/en/latest/)

# Quickstart
(Until i have time to document in more detail)
1. Git clone repo
2. install requirements `pip3 install -r requirements.txt`
3. Generate secret_key `python3 generate_secret_key.py`
4. Edit automarkwatched/settings.py
   * Add generated key to SECRET_KEY = '<here>'
   * Add ip addres or fqdn to ALLOWED_HOSTS
5. Create database `python3 manage.py migrate`
6. Run server `python3 manage.py runserver 0:8000`

Note: If you know how to use gunicorn and nginx, i have configs for both in the root dir. Will add to docs eventually 

# Cronjobs
To automatically run tasks currently you need to create cronjobs that run management commands (until I create a scheduler)  
Example:  
```
# m h  dom mon dow   command
0 */12 * * * /usr/bin/python3 /opt/amw/manage.py filltable
5 * * * * /usr/bin/python3 /opt/amw/manage.py markwatched
30 0 1 * * /usr/bin/python3 /opt/amw/manage.py synctvdb
```

# Upgrades
If a database field has been added you need to migrate the database to use the new features. Its as simple as:  
`python3 manage.py migrate`  
then restart the service

# Pictures
![Home Page](https://i.imgur.com/OO8RFnr.png "Home Page")
![Bulk Edit](https://i.imgur.com/14FYTC7.png "Bulk Edit")
![Show Detail](https://i.imgur.com/DTreE57.png "Show Detail")
![Settings](https://i.imgur.com/eBphkjw.png "Settings")
