(function ($) {
    "use strict";
    $(".seriesToggle").on("change", onScrobbleToggle)
})(jQuery);

function onScrobbleToggle() {
    let scrobble = false
    let id = $(this).attr('id')
    if ($(`#${id}:checked`).length === 1) {
        scrobble = true
    }
    let token = Cookies.get("X-Plex-Token")
    let plexId = Cookies.get("X-Plex-ID")
    let opts = {
        type: "PATCH",
        url: `http://localhost:5309/api/v1/series?scrobble=${scrobble}&ratingKey=${id}`,
        contentType: "application/json",
        headers: {"X-Plex-Token": token, "X-Plex-ID": plexId}
    }
    return $.ajax(opts)
}
