from django.core.management.base import BaseCommand, CommandError
from amw.utilities import plex


class Command(BaseCommand):
    help = 'Marks all episodes in a TV Show set to "Silence" as watched'

    def handle(self, *args, **options):
        server = plex.Plex()
        try:
            server.mark_watched()
            print(self.style.SUCCESS('Successfully set all silenced shows as watched'))
        except:
            CommandError('General Error')

