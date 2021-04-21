(function ($) {
    "use strict";
    $(".seriesToggle").on("change", onScrobbleToggle)
    $("#logout").on("click", () => {
        Cookies.remove('X-Plex-Token')
        Cookies.remove('X-Plex-ID')
        window.location.href = "/login"
    })
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
        url: `${window.location.protocol}//${window.location.hostname}:5309/api/v1/series?scrobble=${scrobble}&ratingKey=${id}`,
        contentType: "application/json",
        headers: {"X-Plex-Token": token, "X-Plex-ID": plexId}
    }
    return $.ajax(opts)
}
