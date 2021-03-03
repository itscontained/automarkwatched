from django.apps import AppConfig


class AmwConfig(AppConfig):
    name = 'amw'
    verbose_name = 'AutoMarkWatched'

    def ready(self):
        from django.db.utils import OperationalError
        from python.amw.utilities import plex
        from python.amw.utilities import thetvdb
        from python.amw import Scheduler
        from django.contrib.auth.models import User

        scheduler = Scheduler()
        try:
            users = User.objects.all()
            for user in users:
                if user.plexuser.current_server:
                    p = plex.Plex(user.username, user.plexuser.authenticationToken)
                    p.connect(user.plexuser.current_server.name)
                    populate_cron = user.scheduling.populate.split()
                    scheduler.scheduler.add_job(
                        p.rectify_show_list,
                        name=f"Populate show list for {user.username}",
                        trigger='cron',
                        minute=populate_cron[0],
                        hour=populate_cron[1],
                        day=populate_cron[2],
                        month=populate_cron[3],
                        day_of_week=populate_cron[4],
                        misfire_grace_time=15
                    )
                    sync_cron = user.scheduling.sync.split()
                    tvdb = thetvdb.TheTVDB(user.username)
                    scheduler.scheduler.add_job(
                        tvdb.sync_shows,
                        name=f"Sync show list for {user.username}",
                        trigger='cron',
                        minute=sync_cron[0],
                        hour=sync_cron[1],
                        day=sync_cron[2],
                        month=sync_cron[3],
                        day_of_week=sync_cron[4],
                        misfire_grace_time=15
                    )
                    mark_watched_cron = user.scheduling.mark_watched.split()
                    scheduler.scheduler.add_job(
                        p.mark_watched,
                        name=f"Mark Watched for {user.username}",
                        trigger='cron',
                        minute=mark_watched_cron[0],
                        hour=mark_watched_cron[1],
                        day=mark_watched_cron[2],
                        month=mark_watched_cron[3],
                        day_of_week=mark_watched_cron[4],
                        misfire_grace_time=15
                    )
        except OperationalError:
            pass
