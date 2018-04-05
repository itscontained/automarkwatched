from django.urls import path
from django.contrib import admin
from . import views

urlpatterns = [
    path('', views.HomeView.as_view()),
    path('settings', views.SettingsView.as_view()),
    path('settings/filltable/', views.filltable),
    path('settings/synctvdb/', views.syncTVDB),
    path('settings/markwatched/', views.markWatched),
    path('<int:show_pk>/', views.ShowDetailView.as_view()),
]

