let TVShowTable = {};

window.addEventListener('DOMContentLoaded', event => {
    const tvShowTable = document.getElementById('tvShowsTable');
    if (tvShowTable) {
        TVShowTable = new simpleDatatables.DataTable(tvShowTable, {
            columns: [
                {
                    select: 0,
                    sortable: false,
                    render: function(data, cell, row) {
                        let d = JSON.parse(data)
                        let checked = d.scrobble ? "checked" : ""
                        return "<mwc-switch class='series-toggle' data-rating-key='" + d.ratingKey + "' " + checked + "></mwc-switch>"
                    }
                },
                { select: 1, hidden: false},
                { select: 2, sort: "asc"}
            ]
        });
        TVShowTable.on('datatable.init', function () {
            $.when(
                apiCall("servers"),
                apiCall("libraries"),
                apiCall("series")
            ).done(function (p1, p2, p3) {
                let servers = p1[0]
                let libraries = p2[0]
                let series = p3[0]
                let data = []
                for (const [ratingKey, s] of Object.entries(series)) {
                    data.push({
                        Scrobble: JSON.stringify({
                            scrobble: s.scrobble,
                            ratingKey: ratingKey,
                        }),
                        "RatingKey": ratingKey,
                        Title: s.title,
                        Year: s.year.toString(),
                        Library: libraries[s.library_uuid].title,
                        Server: servers[s.server_machine_identifier].name
                    })
                }
                TVShowTable.insert(data)
            })
        })
        $("table#tvShowsTable").on("click", "mwc-switch", onScrobbleToggle)
    }
});
