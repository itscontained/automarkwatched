from django.views import View
from django.contrib import messages
from django.shortcuts import render
from django.views.generic import TemplateView, ListView
from django.http import HttpResponseRedirect, HttpResponse
from django.contrib.auth.decorators import login_required
from django.contrib.auth.models import User
from django.contrib.auth.mixins import LoginRequiredMixin

from .utilities import plex, thetvdb
from .models import TVShow, PlexMediaServer
from .forms import ServerForm


class HomeView(LoginRequiredMixin, ListView):
    template_name = "amw/home.html"
    context_object_name = "tv_shows"

    def get_queryset(self):
        context = TVShow.objects.order_by('title')
        return context

    def post(self, request, *args, **kwargs):
        bulkeditinfo = [ (show, request.POST[show]) for show in request.POST if 'csrfmiddlewaretoken' not in show ]
        valid = {
            'True': True,
            'False': False
        }
        for pkid, choice in bulkeditinfo:
            show = TVShow.objects.filter(id=pkid)[0]
            if show.silenced != valid[choice]:
                show.silenced = valid[choice]
                show.save()
                if valid[choice]:
                    print('Silenced {}'.format(show.title))
                else:
                    print('Unsilenced {}'.format(show.title))
        return HttpResponseRedirect('/')


class SettingsView(LoginRequiredMixin, TemplateView):
    form_class = ServerForm
    initial = {'key': 'value'}
    template_name = "amw/settings.html"

    def get(self, request, *args, **kwargs):
        plexuser = request.user.plexuser

        form = self.form_class()
        form.fields['servers'].choices = tuple(plexuser.servers.all().values_list())
        form.initial['servers'] = plexuser.servers.all().filter(name=plexuser.current_server.name).values_list()[0][0]
        context = {
            'form': form
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        print()
        form = ServerForm(request.POST)
        selected_server = PlexMediaServer.objects.get(pk=request.POST['servers'])
        if request.user.plexuser.current_server != selected_server:
            request.user.plexuser.current_server = selected_server
            request.user.save()
            messages.success(request, "Active PlexMediaServer successfully changed")
        else:
            print('shit')
        return HttpResponseRedirect('/settings')


class ShowDetailView(TemplateView):
    template_name = "amw/showdetail.html"

    def get(self, request, show_pk, *args, **kwargs):
        context = { 'show': TVShow.objects.get(pk=show_pk) }
        return render(request, self.template_name, context)

    def post(self, request, show_pk, *args, **kwargs):
        show = TVShow.objects.get(pk=show_pk)

        if show.silenced:
            show.silenced = False
            print('Unsilenced {}'.format(show.title))
        else:
            show.silenced = True
            print('Silenced {}'.format(show.title))

        show.save()

        return HttpResponseRedirect('/{}'.format(show_pk))


def filltable(request):

    if request.method == "POST":
        u = request.user.plexuser
        p = plex.Plex(u.user, u.authenticationToken)
        p.connect(u.current_server.name)
        p.rectify_show_list()
        messages.success(request, 'Success! TV Show table populated!')

        return HttpResponseRedirect('/')


def syncTVDB(request):

    if request.method == "POST":
        server = thetvdb.TheTVDB()
        server.syncShows()
        messages.success(request, 'Success! TV show continuing status synced with TheTVDB')

        return HttpResponseRedirect('/')


def markWatched(request):

    if request.method == "POST":
        server = plex.Plex()
        server.mark_watched()
        messages.success(request, 'Success! TV shows have been marked watched')

        return HttpResponseRedirect('/settings')

