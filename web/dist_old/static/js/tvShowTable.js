window.addEventListener('DOMContentLoaded', event => {
    // Simple-DataTables
    // https://github.com/fiduswriter/Simple-DataTables/wiki

    const tvShowTable = document.getElementById('tvShowTable');
    if (tvShowTable) {
        new simpleDatatables.DataTable(tvShowTable, {
            columns: [
                {
                    select: 0,
                    sortable: false,
                    render: function(data, cell, row) {
                        let checked = data === "true" ? "checked" : ""
                        return "<mwc-switch "+checked+"></mwc-switch>"
                    }
                },
                { select: 1, hidden: true},
                { select: 2, sort: "asc"}
            ]
        });
    }
});
