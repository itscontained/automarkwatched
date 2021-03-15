const plexLoginButton = $("#PlexLoginButton");
const plexLoadServersButton = $("#PlexLoadServersButton");
const plexLoadLibrariesButton = $("#PlexLoadLibrariesButton");
const plexLoadSeriesButton = $("#PlexLoadSeriesButton");
let pollSSOWindowID;

(function ($) {
    "use strict";

    // load plex servers
    plexLoginButton.on("click", plexLogin)
    plexLoadServersButton.on("click", plexLoadServers)
    plexLoadLibrariesButton.on("click", plexLoadLibraries)
    plexLoadSeriesButton.on("click", plexSyncSeries)
    $("#wizard1-next").on("click", () => {
        $("#wizard1-tab").removeClass("active")
        $("#wizard2-tab").addClass("active")
    })
    $("#wizard2-next").on("click", () => {
        setPlexServers()
        $("#wizard2-tab").removeClass("active")
        $("#wizard3-tab").addClass("active")
    })
    $("#wizard3-next").on("click", () => {
        plexSetLibraries()
        $("#wizard3-tab").removeClass("active")
        $("#wizard4-tab").addClass("active")
    })
    $("#wizard4-finish").on("click", () => {
        window.location.href = "/"
    })
})(jQuery);

function apiCall(endpoint, data) {
    let token = Cookies.get("X-Plex-Token")
    let opts = {
        type: "GET",
        url: `http://localhost:5309/api/v1/${endpoint}`,
        contentType: "application/json",
        headers: {"X-Plex-Token": token}
    }
    if (data !== undefined) {
        opts.type = "POST"
        opts.data = JSON.stringify(data)
    }
    return $.ajax(opts)
}
function plexLogin() {
    plexLoginButton.attr("disabled", true)
    plexLoginButton.children().toggle()
    let loginData = null
    if (Cookies.get("X-Plex-Token") !== undefined) {
        getUser()
        return
    }
    $.get( "http://localhost:5309/api/v1/login").done((data) => {
        loginData = data
        let features = "resizable,scrollbars,status,width=700,height=600"
        let ssoWindow = window.open(data.url, "Plex.TV SSO", features)
        if (ssoWindow != null) {
            pollSSOWindowID = setInterval(pollSSOWindow, 1000, ssoWindow, loginData)
        }
    });
}

function plexLoadServers() {
    plexLoadServersButton.attr("disabled", true)
    plexLoadServersButton.children().toggle()
    apiCall("servers").done((data) => {
        for (const [machineID, server] of Object.entries(data)) {
            let t = $("#serverItemTemplateInput").clone(true, true)
            t.find("input").data("object", server)
                .attr("id", server.Name)
                .attr("checked", true)
            t.find("label").attr("for", server.Name)
            t.find("h1").text(server.Name)
            t.find("small.url").text(`${server.Scheme}://${server.Host}:${server.Port}`)
            t.find("small.mid").text(machineID)
            t.find("small.version").text(server.Version)
            t.removeClass("d-none")
            $("#serverItemTemplateInput").parent().append(t)
        }
        plexLoadServersButton.children().toggle()
        plexLoadServersButton.removeAttr("disabled")
        $("#wizard2").find("button.btn-primary").removeAttr("disabled")
    })
}

function setPlexServers() {
    let servers = Object();
    $("#wizard2").find("input:checked").each(function () {
        let server = $(this).data("object")
        servers[server.MachineIdentifier] = server
    })
    apiCall("servers", servers).done(() => {

    })
}
function pollSSOWindow(ssoWindow, loginData) {
    if (ssoWindow == null || ssoWindow.closed) {
        console.log("sso window closed")
        clearInterval(pollSSOWindowID)
    }
    let pinData = {
        code: loginData.pin.code,
        'X-Plex-Client-Identifier': loginData.pin.clientIdentifier,
    }
    $.get(`https://plex.tv/api/v2/pins/${loginData.pin.id}.json`, pinData).done((data) => {
        if (data.authToken !== null) {
            clearInterval(pollSSOWindowID)
            Cookies.set("X-Plex-Token", data.authToken)
            getUser()
        }
    })
}

function getUser() {
    let token = Cookies.get("X-Plex-Token")
    $.get("http://localhost:5309/api/v1/user", {"X-Plex-Token": token}).done((data) => {
        $("#wizard1-content").children("div.row-cols-1.collapse").collapse()
        $("#wizard1-content").children("div.row.collapse").show()
        $("#wizard1-content").find("button.btn-primary").removeAttr("disabled")
        $("#wizard1-content").find("img").attr("src", data.thumb)
        $("#wizard1-content").find("h1").text(data.username)
        if (!data.exists) {
            console.log(data)
            data.owner = true
            $.ajax({
                type: "POST",
                url: `http://localhost:5309/api/v1/user`,
                data: JSON.stringify(data),
                contentType: "application/json",
                headers: {"X-Plex-Token": token}
            })
        }
    })
}

function plexLoadLibraries() {
    plexLoadLibrariesButton.attr("disabled", true)
    plexLoadLibrariesButton.children().toggle()
    apiCall("libraries").done((data) => {
        for (const [uuid, library] of Object.entries(data)) {
            let t = $("#libraryItemTemplateInput").clone(true, true)
            t.find("input").data("object", library)
                .attr("id", library.Name)
                .attr("checked", true)
            t.find("label").attr("for", library.Name)
            t.find("h1").text(library.title)
            t.find("small.server-name").text(uuid)
            t.find("small.agent").text(library.agent)
            t.find("small.scanner").text(library.scanner)
            t.removeClass("d-none")
            $("#libraryItemTemplateInput").parent().append(t)
        }
        plexLoadLibrariesButton.children().toggle()
        plexLoadLibrariesButton.removeAttr("disabled")
        $("#wizard3").find("button.btn-primary").removeAttr("disabled")
    })
}

function plexSetLibraries() {
    let libraries = Object();
    $("#wizard3").find("input:checked").each(function () {
        let library = $(this).data("object")
        libraries[library.uuid] = library
    })
    apiCall("libraries", libraries).done(() => {

    })
}

function plexSyncSeries() {
    plexLoadSeriesButton.attr("disabled", true)
    plexLoadSeriesButton.children().toggle()
    apiCall("series/sync").done(() => {
        plexLoadSeriesButton.children().toggle()
        $("#wizard4").find("button.btn-primary").removeAttr("disabled")
    })
}
