from os import environ as env
from django.conf import settings
from django.db import migrations, models
import django.db.models.deletion


def generate_superuser(apps, schema_editor):
    from django.contrib.auth.models import User

    django_su_name = env.get('AMW_SUPERUSER_USERNAME')
    django_su_email = env.get('AMW_SUPERUSER_EMAIL')
    django_su_password = env.get('AMW_SUPERUSER_PASSWORD')

    superuser = User.objects.create_superuser(
        username=django_su_name,
        email=django_su_email,
        password=django_su_password)

    superuser.save()


class Migration(migrations.Migration):

    initial = True

    dependencies = [
        migrations.swappable_dependency(settings.AUTH_USER_MODEL),
    ]

    operations = [
        migrations.CreateModel(
            name='PlexMediaServer',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=200, unique=True)),
            ],
        ),
        migrations.CreateModel(
            name='TVShow',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('title', models.CharField(max_length=200)),
                ('silenced', models.BooleanField()),
                ('continuing', models.BooleanField()),
                ('tvdbid', models.IntegerField()),
                ('banner_url', models.CharField(default='', max_length=20)),
                ('user', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='Scheduling',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('populate', models.CharField(default='0 2 * * *', max_length=200)),
                ('sync', models.CharField(default='0 4 * * 6', max_length=200)),
                ('mark_watched', models.CharField(default='*/15 * * * *', max_length=200)),
                ('user', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE,
                                              to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.CreateModel(
            name='PlexUser',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('authenticationToken', models.CharField(blank=True, max_length=200)),
                ('current_server', models.ForeignKey(null=True, on_delete=django.db.models.deletion.SET_NULL,
                                                     related_name='current_users', to='amw.PlexMediaServer')),
                ('servers', models.ManyToManyField(blank=True, related_name='associated_users',
                                                   to='amw.PlexMediaServer')),
                ('user', models.OneToOneField(on_delete=django.db.models.deletion.CASCADE,
                                              to=settings.AUTH_USER_MODEL)),
            ],
        ),
        migrations.RunPython(generate_superuser),
    ]
