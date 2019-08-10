from apscheduler.schedulers.background import BackgroundScheduler


def start():
    scheduler = BackgroundScheduler()
    scheduler.add_job('somefunctionhere', 'interval', minutes=5)
    scheduler.start()
