"use strict";
// interface FindResponse {
// 	running: App[];
// 	deployed: App[];
// }
var baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);
var token = document.cookie.split("; ").filter(function (e) { return e.startsWith("Authorization"); })[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
var tokenData = jwt_decode(token);
function init() {
    var url = baseUrl;
    url.pathname = "/api/find";
    fetch(url.href).then(function (j) {
        j.json().then(function (res) {
            console.log(res);
            var appsD = res.deployed != null ? res.deployed : [];
            var apps = res.running != null ? res.running : [];
            appsD.forEach(function (a) {
                if (apps.filter(function (app) { return a.id == app.id; }).length == 1) {
                    document.querySelector("#appContainer").innerHTML += appTemplate(apps.find(function (app) { return a.id == app.id; }), true);
                }
                else {
                    document.querySelector("#appContainer").innerHTML += appTemplate(a, false);
                }
            });
        }).catch(function (err) { return console.log(err); });
    }).catch(function (err) { return console.log(err); });
}
function getButton(id, action) {
    var icon, color, text;
    switch (action) {
        case "run":
            color = "success";
            icon = "running";
            text = "Run";
            break;
        case "kill":
            color = "warning";
            icon = "skull";
            text = "Kill";
            break;
        case "update":
            color = "info";
            icon = "sync";
            text = "Update";
            break;
        case "remove":
            color = "danger";
            icon = "trash";
            text = "Remove";
            break;
    }
    return "<button class=\"btn btn-" + color + "\" data-action=\"" + action + "\" data-id=\"" + id + "\" onclick=\"doAction(event)\"><i class=\"fas fa-" + icon + " fa-2x\"></i><br>" + text + "</button>";
}
function dateTemplate(dateString) {
    return new Date(dateString).toLocaleString();
}
function runnerIcon(runner) {
    var r = "";
    switch (runner) {
        case "node":
            r = "node";
            break;
        case "web":
            r = "html5";
            break;
    }
    return "<i class=\"fab fa-" + r + " fa-2x\"></i>";
}
function appTemplate(app, running) {
    return " <div class=\"card\">\n            <div class=\"card-header\" id=\"heading" + app.id + "\">\n\t\t\t\t<span class=\"float-right " + (running ? "online text-success" : "offline text-danger") + "\">" + (running ? "Online <i class=\"fas fa-globe\"></i>" : "Offline <i class=\"fas fa-times-circle\"></i>") + "</span>\n                <h3 class=\"mb-2\" style=\"cursor: pointer;\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.name + "\n                </h3>\n                <h6 class=\"mb-0 text-muted\">\n\t\t\t\t\t" + app.repo + "\n\t\t\t\t</h6>\n            </div>\n            <div id=\"collapse" + app.id + "\" class=\"collapse\" aria-labelledby=\"heading" + app.id + "\" data-parent=\"#appContainer\">\n                <div class=\"card-body row p-0\">\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>ID:</span><span>" + app.id + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Name:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Repo:</span><span>" + app.repo + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Root:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Port:</span><span>" + (app.port == 0 ? "none" : app.port) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Hostname:</span><span>" + app.hostname + "</span>\n                        </li>\n                    </ul>\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Deployed:</span><span>" + dateTemplate(app.deployed) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastUpdated:</span><span>" + dateTemplate(app.last_updated) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastRun:</span><span>" + dateTemplate(app.last_run) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Uptime:</span><span>" + (running ? app.uptime.replace(/\.(.*?)s/g, "s") : "offline") + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Runner:</span><span>" + runnerIcon(app.runner) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Pid:</span><span>" + (app.pid == 0 ? "offline" : app.pid) + "</span>\n                        </li>\n                    </ul>\n                </div>\n                <div class=\"card-footer text-right\">\n                \t" + (running ? getButton(app.id, "kill") : getButton(app.id, "run")) + "\n    \t\t\t\t" + getButton(app.id, "update") + "\n    \t\t\t\t" + getButton(app.id, "remove") + "\n\t\t\t\t</div>\n            </div>\n        </div>";
}
function doAction(event) {
    var url = baseUrl;
    var btn = event.target;
    var action = btn.attributes.getNamedItem("data-action").value;
    var id = btn.attributes.getNamedItem("data-id").value;
    var data = { app: id };
    url.pathname = "/api/" + action;
    fetch(url.href, {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    })
        .then(function (res) { return res.status == 200 ? location.reload() : null; })
        .catch(function (err) { return console.log(err); });
}
init();
