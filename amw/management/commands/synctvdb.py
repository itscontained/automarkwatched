from django.core.management.base import BaseCommand, CommandError
from amw.utilities import thetvdb


class Command(BaseCommand):
    help = 'Marks all episodes in a TV Show set to "Silence" as watched'

    def handle(self, *args, **options):
        server = thetvdb.TheTVDB()
        try:
            server.sync_shows()
            print(self.style.SUCCESS('Successfully synced all shows with TheTVDB'))
        except:
            CommandError('General Error')

