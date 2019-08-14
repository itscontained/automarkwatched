from os import environ as env
from django.db import migrations


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

    dependencies = [
        ('amw', '0001_initial'),
    ]

    operations = [
        migrations.RunPython(generate_superuser),
    ]
