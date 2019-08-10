from django.apps import AppConfig


class AmwConfig(AppConfig):
    name = 'amw'
    verbose_name = 'AutoMarkWatched'

    def ready(self):
        from amw import scheduler
        scheduler.start()
