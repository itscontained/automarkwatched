from requests import Session

from django.contrib import messages

from ..models import TVShow


class TheTVDB(object):
    baseurl = 'https://api.thetvdb.com'

    def __init__(self, user):
        self.user = user
        self.session = Session()
        payload = {'apikey': 'RSD94YQUFALZO3DZ'}
        token = self.session.post(f"{TheTVDB.baseurl}/login", json=payload).json()['token']
        self.session.headers.update({'Authorization': f'Bearer {token}'})

    def get_show_info(self, tv_show):
        g = self.session.get(f'{TheTVDB.baseurl}/series/{tv_show.tvdbid}')
        return None if not g.ok else g.json()['data']

    def sync_shows(self):
        for show in TVShow.objects.filter(user=self.user):
            statuses = {
                'Continuing': True,
                'Ended': False
            }
            showinfo = self.get_show_info(show)
            if showinfo:
                if showinfo.get('status'):
                    if show.continuing != statuses[showinfo['status']]:
                        show.continuing = statuses[showinfo['status']]
                        show.save()
                        print(f"Updated continuing status for {show.title} to {showinfo['status']}")
                if showinfo['banner']:
                    if show.banner_url != showinfo['banner']:
                        show.banner_url = showinfo['banner']
                        show.save()
                        print(f"updated Banner URL for {show.title}")
            else:
                print(f"error pulling info for {show.title}")
