from django.urls import path
from . import views

urlpatterns = [
    path('', views.HomeView.as_view(), name='home'),
    path('settings', views.SettingsView.as_view()),
    path('scheduler', views.SchedulerView.as_view()),
    path('api/filltable/', views.filltable),
    path('api/synctvdb/', views.sync_tvdb),
    path('api/markwatched/', views.mark_watched),
    path('<int:show_pk>/', views.ShowDetailView.as_view()),
]

