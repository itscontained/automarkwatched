from apscheduler.schedulers.background import BackgroundScheduler

from amw.utilities import plex


class Scheduler(object):
    def __init__(self):
        self.scheduler = BackgroundScheduler()
        self.scheduler.start()
