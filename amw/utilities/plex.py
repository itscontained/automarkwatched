from django.contrib import messages

from plexapi.server import PlexServer

from ..models import TVShow, ServerInfo


class Plex(object):

    def __init__(self):
        self.url = 'http://{}'.format(ServerInfo.objects.all()[0].url)
        self.token = ServerInfo.objects.all()[0].token
        self.server = PlexServer(self.url, self.token)
        self.tv_show_section = self.server.library.section('TV Shows')

    def get_tv_show_list(self):
        tv_show_list = [(show.title, ''.join([s for s in list(show.guid) if s.isdigit()])) for show in self.tv_show_section.all()]
        return tv_show_list

    def rectify_show_list(self):
        show_list = self.get_tv_show_list()
        for show, thetvdbid in show_list:
            if not TVShow.objects.filter(title=show):
                missing_show = TVShow(title=show, silenced=False, continuing=False, tvdbid=thetvdbid)
                missing_show.save()
                print('added missing show {}, with tvdbid {}'.format(show, thetvdbid))

    def mark_watched(self):
        silenced_show_objects = TVShow.objects.filter(silenced=True)
        silenced_show_list = [show.title for show in silenced_show_objects]
        for episode in self.tv_show_section.searchEpisodes(unwatched=True):
            if episode.grandparentTitle in silenced_show_list:
                print('marking ', episode.title, ' from ', episode.grandparentTitle, ' watched')
                episode.markWatched()

