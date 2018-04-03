from django.db import models

class TVShow(models.Model):
    title = models.CharField(max_length=200)
    silenced = models.BooleanField()
    continuing = models.BooleanField()
    tvdbid = models.CharField(max_length=20, default='')

class ServerInfo(models.Model):
    url = models.CharField(max_length=200)
    token = models.CharField(max_length=100)
