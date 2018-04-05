from django.core.management.base import BaseCommand, CommandError
from amw.utilities import plex


class Command(BaseCommand):
    help = 'Fills table with missing TV Shows from Plex'

    def handle(self, *args, **options):
        server = plex.Plex()
        try:
            server.rectify_show_list()
            print(self.style.SUCCESS('Successfully updated table'))
        except:
            CommandError('General Error')

