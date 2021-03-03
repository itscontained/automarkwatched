from apscheduler.schedulers.background import BackgroundScheduler


class Scheduler(object):
    def __init__(self):
        self.scheduler = BackgroundScheduler()
        self.scheduler.start()
