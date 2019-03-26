"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function appTemplate(app, running) {
    return " <div class=\"card\">\n            <div class=\"card-header\" id=\"heading" + app.id + "\">\n\t\t\t\t<span class=\"float-right " + (running ? "online text-success" : "offline text-danger") + "\">" + (running ? "Online &bull;" : "Offline") + "</span>\n                <h2 class=\"mb-0\" style=\"cursor: pointer;\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.name + "\n                </h2>\n                <h5 class=\"mb-0 text-muted\">\n\t\t\t\t\t" + app.repo + "\n\t\t\t\t</h5>\n            </div>\n            <div id=\"collapse" + app.id + "\" class=\"collapse\" aria-labelledby=\"heading" + app.id + "\" data-parent=\"#appContainer\">\n                <div class=\"card-body row\">\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>ID:</span><span>" + app.id + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Name:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Repo:</span><span>" + app.repo + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Root:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Port:</span><span>" + app.port + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Hostname:</span><span>" + app.hostname + "</span>\n                        </li>\n                    </ul>\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Deployed:</span><span>" + app.deployed + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastUpdated:</span><span>" + app.last_updated + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastRun:</span><span>" + app.last_run + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Uptime:</span><span>" + (running ? app.uptime : "offline") + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Runner:</span><span>" + app.runner + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Pid:</span><span>" + app.pid + "</span>\n                        </li>\n                    </ul>\n                </div>\n            </div>\n        </div>";
}
;
var url = new URL(window.location.href);
var token = document.cookie.split("; ").filter(function (e) { return e.startsWith("Authorization"); })[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
var tokenData = jwt_decode(token);
console.log(tokenData);
fetch(url.href + "/api/find").then(function (j) {
    j.json().then(function (res) {
        console.log(res);
        var appsD = res.deployed;
        var apps = res.running != null ? res.running : [];
        appsD.forEach(function (a) {
            var id = a.id;
            if (apps.filter(function (app) { return id == app.id; }).length == 1) {
                document.querySelector("#appContainer").innerHTML += appTemplate(apps.find(function (app) { return id == app.id; }), true);
            }
            else {
                document.querySelector("#appContainer").innerHTML += appTemplate(a, false);
            }
        });
    }).catch(function (err) { return console.log(err); });
}).catch(function (err) { return console.log(err); });
