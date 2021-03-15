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
    let opts = {
        type: "PATCH",
        url: `http://localhost:5309/api/v1/series/scrobble/${id}?scrobble=${scrobble}`,
        contentType: "application/json",
        headers: {"X-Plex-Token": token}
    }
    return $.ajax(opts)
}
