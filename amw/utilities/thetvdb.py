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
        if g.status_code != 200:
            showinfo = ''
        else:
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
                bannerurl = showinfo['banner']
                if tvdbshowstatus:
                    if dbshowstatus != statuses[tvdbshowstatus]:
                        show.continuing = statuses[tvdbshowstatus]
                        show.save()
                        print('updated continuing status for {} to {}'.format(show.title, tvdbshowstatus))
                if bannerurl:
                    if bannerurl != show.bannerurl:
                        show.bannerurl = bannerurl
                        print('updated Banner URL for {}'format(show.title))
            else:
                print('error pulling info for {}'.format(show.title))
