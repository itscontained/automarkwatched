import requests

from django.contrib import messages

from ..models import TVShow

class TheTVDB(object):

    def __init__(self):
        amwapikey = {'apikey': '0E4FB4A79D7C3D09'}
        self.apiurl = 'https://api.thetvdb.com'
        self.token = requests.post('{}/login'.format(self.apiurl), json=amwapikey).json()['token']
        self.headers = {'Authorization': 'Bearer {}'.format(self.token)}

    def getShowInfo(self, tvshow):
        g = requests.get('{}//series/{}'.format(self.apiurl, tvshow.tvdbid), headers=self.headers)
        showinfo = g.json()['data']
        return showinfo

    def syncShows(self):
        for show in TVShow.objects.all():
            statuses = {
                'Continuing': True,
                'Ended': False
            }
            showinfo = self.getShowInfo(show)
            if showinfo:
                dbshowstatus = show.continuing
                tvdbshowstatus = showinfo['status']
                if dbshowstatus != statuses[tvdbshowstatus]:
                    show.continuing = statuses[tvdbshowstatus]
                    show.save()
                    print('updated continuing status for {} to {}'.format(show.title, tvdbshowstatus))

        messages.success(request, 'Success! TV show continuing status synced with TheTVDB')
