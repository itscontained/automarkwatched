FROM amd64/python:3.7.4-alpine

LABEL maintainers="dirtycajunrice"

WORKDIR /app

COPY /python/requirements.txt /manage.py /start.sh /app/

COPY /python/amw /app/amw

COPY /python/automarkwatched /app/automarkwatched

RUN apk add --no-cache tzdata && \
    pip install --no-cache-dir -r /app/requirements.txt

VOLUME /data
VOLUME /static_files

ENTRYPOINT ["./start.sh"]
