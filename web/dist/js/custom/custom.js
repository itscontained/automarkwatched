(function ($) {
    "use strict";
    $(".series-toggle").on("change", onScrobbleToggle)
    $("#logout").on("click", () => {
        Cookies.remove('X-Plex-Token')
        Cookies.remove('X-Plex-ID')
        window.location.href = "/login"
    })
})(jQuery);


function onScrobbleToggle() {
    let mwc = $(this)
    let scrobble = !mwc[0].__checked
    let ratingKey = mwc.data("rating-key")
    let token = Cookies.get("X-Plex-Token")
    let plexId = Cookies.get("X-Plex-ID")
    let opts = {
        type: "PATCH",
        url: `${window.location.protocol}//${window.location.hostname}:5309/api/v1/series?scrobble=${scrobble}&ratingKey=${ratingKey}`,
        contentType: "application/json",
        headers: {"X-Plex-Token": token, "X-Plex-ID": plexId}
    }
    return $.ajax(opts)
}
