from apscheduler.schedulers.background import BackgroundScheduler

from amw.utilities import plex


def start():
    #server = plex.Plex()
    scheduler = BackgroundScheduler()
    #scheduler.add_job(server.rectify_show_list, 'interval', minutes=30)
    scheduler.start()
