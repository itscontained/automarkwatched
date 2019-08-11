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
from .forms import ServerForm, SchedulerForm


class HomeView(LoginRequiredMixin, ListView):
    template_name = "amw/home.html"
    context_object_name = "tv_shows"

    def get_queryset(self):
        context = TVShow.objects.filter(user=self.request.user).order_by('title')
        return context

    def post(self, request, *args, **kwargs):
        bulkeditinfo = [(show, request.POST[show]) for show in request.POST if 'csrfmiddlewaretoken' not in show]
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
        form.fields['servers'].choices = list(plexuser.servers.all().values_list())
        if plexuser.current_server:
            form.initial['servers'] = plexuser.servers.all().filter(
                name=plexuser.current_server.name).values_list()[0][0]
        else:
            form.fields['servers'].choices.insert(0, (0, "Choose..."))
            form.initial['servers'] = 0
        context = {
            'form': form,
            'tv_shows': TVShow.objects.all()
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        selected_server = PlexMediaServer.objects.get(pk=request.POST['servers'])
        if request.user.plexuser.current_server != selected_server:
            request.user.plexuser.current_server = selected_server
            request.user.save()
            messages.success(request, "Active PlexMediaServer successfully changed")
        else:
            print('shit')
        return HttpResponseRedirect('/settings')


class ShowDetailView(LoginRequiredMixin, TemplateView):
    template_name = "amw/showdetail.html"

    def get(self, request, show_pk, *args, **kwargs):
        context = {'show': TVShow.objects.get(pk=show_pk)}
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


class SchedulerView(LoginRequiredMixin, TemplateView):
    form_class = SchedulerForm
    initial = {'key': 'value'}
    template_name = "amw/scheduler.html"

    def get(self, request, *args, **kwargs):
        scheduling = request.user.scheduling
        self.initial = {
            'cron_populate': scheduling.populate,
            'cron_sync': scheduling.sync,
            'cron_mark_watched': scheduling.mark_watched
        }
        form = self.form_class(initial=self.initial)
        context = {
            'form': form
        }
        return render(request, self.template_name, context)

    def post(self, request, *args, **kwargs):
        if 'cron_populate' in request.POST:
            request.user.scheduling.populate = request.POST['cron_populate']
            messages.success(request, "Updated scheduling. Please restart")
            request.user.save()
        elif 'cron_sync' in request.POST:
            request.user.scheduling.sync = request.POST['cron_sync']
            messages.success(request, "Updated scheduling. Please restart")
            request.user.save()
        elif 'cron_mark_watched' in request.POST:
            request.user.scheduling.mark_watched = request.POST['cron_mark_watched']
            messages.success(request, "Updated scheduling. Please restart")
            request.user.save()
        return HttpResponseRedirect('/scheduler')


def filltable(request):
    if request.method == "POST":
        u = request.user.plexuser
        p = plex.Plex(u.user, u.authenticationToken)
        p.connect(u.current_server.name)
        p.rectify_show_list()
        messages.success(request, 'Success! TV Show table populated!')

        return HttpResponseRedirect('/')


def sync_tvdb(request):
    if request.method == "POST":
        server = thetvdb.TheTVDB(request.user)
        server.sync_shows()
        messages.success(request, 'Success! TV show continuing status synced with TheTVDB')

        return HttpResponseRedirect('/')


def mark_watched(request):
    if request.method == "POST":
        server = plex.Plex(request.user.username, request.user.plexuser.authenticationToken)
        server.mark_watched()
        messages.success(request, 'Success! TV shows have been marked watched')

        return HttpResponseRedirect('/settings')

