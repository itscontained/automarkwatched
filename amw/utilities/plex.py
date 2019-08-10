from django.contrib import messages

from amw.models import TVShow
from plexapi.myplex import MyPlexAccount


class Plex(object):

    def __init__(self, username, token):
        self.account = MyPlexAccount(username=username, token=token)
        self.server = None
        self.tv_show_section = None

    def connect(self, resource):
        self.server = self.account.resource(resource).connect()
        self.tv_show_section = self.server.library.section('TV Shows')

    def get_tv_show_list(self):
        return [(show.title,int(''.join([s for s in show.guid if s.isdigit()]))) for show in self.tv_show_section.all()]

    def rectify_show_list(self):
        show_list = self.get_tv_show_list()
        for show, tvdbid in show_list:
            if not TVShow.objects.filter(tvdbid=tvdbid):
                missing_show = TVShow(title=show, silenced=False, continuing=False, tvdbid=tvdbid)
                missing_show.save()
                print('added missing show {}, with tvdbid {}'.format(show, tvdbid))

    def mark_watched(self):
        silenced_show_objects = TVShow.objects.filter(silenced=True)
        silenced_show_list = [show.title for show in silenced_show_objects]
        for episode in self.tv_show_section.searchEpisodes(unwatched=True):
            if episode.grandparentTitle in silenced_show_list:
                print('marking ', episode.title, ' from ', episode.grandparentTitle, ' watched')
                episode.markWatched()

