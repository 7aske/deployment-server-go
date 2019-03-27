"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var Store = /** @class */ (function () {
    function Store(initialState) {
        var _this = this;
        this._state = {};
        this.state = {};
        Object.keys(initialState).forEach(function (key) {
            _this._state[key] = { value: initialState[key], actions: [] };
        });
        this.state = initialState;
    }
    Store.prototype.setState = function (name, state) {
        if (Object.keys(this._state).indexOf(name) == -1) {
            throw new Error("State must be registered first");
        }
        else {
            this.set(name, state);
            if (this._state[name].actions) {
                this._state[name].actions.forEach(function (action) {
                    action();
                });
            }
        }
        return this.state[name];
    };
    Store.prototype.registerState = function (name, initialState) {
        if (Object.keys(this.state).indexOf(name) != -1) {
            throw new Error("State already exists");
        }
        else {
            this.set(name, initialState);
        }
    };
    Store.prototype.getState = function (name) {
        if (Object.keys(this.state).indexOf(name) == -1) {
            throw new Error("State is not registered - '" + name + "'");
        }
        else {
            return this.state[name];
        }
    };
    Store.prototype.subscribe = function (name, actions) {
        if (Object.keys(this._state).indexOf(name) == -1) {
            throw new Error("State is not registered");
        }
        else {
            this._state[name].actions = actions;
        }
    };
    Store.prototype.getStateObject = function () {
        return this._state;
    };
    Store.prototype.set = function (name, state) {
        if (Object.keys(this.state).indexOf(name) == -1) {
            this.state[name] = state;
            this._state[name] = { value: state, actions: [] };
        }
        else {
            this._state[name].value = state;
            this.state[name] = state;
        }
    };
    return Store;
}());
var initialState = {
    isModalUp: false,
    loading: false,
};
var store = new Store(initialState);
store.subscribe("isModalUp", [updateModal]);
store.subscribe("loading", [toggleLoader]);
var baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);
var token = document.cookie.split("; ").filter(function (e) { return e.startsWith("Authorization"); })[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
var tokenData = jwt_decode(token);
var appContainer = document.querySelector("#appContainer");
var modal = document.querySelector("#modalDialog");
var modalForm = document.querySelector("#modalDialog form");
var modalConfirm = document.querySelector("#btnModalConfirm");
modalConfirm.addEventListener("click", function (e) { return doForm(e); });
var modalCancel = document.querySelector("#btnModalCancel");
var searchInp = document.querySelector("#searchInp");
searchInp.addEventListener("keydown", function () {
    updateApps(searchInp.value);
});
var searchBtn = document.querySelector("#searchBtn");
searchBtn.addEventListener("click", function (e) {
    updateApps(searchInp.value);
});
var deployBtn = document.querySelector("#deployBtn");
deployBtn.addEventListener("click", function (e) {
});
$("#modalDialog")
    .on("shown.bs.modal", function () { return store.setState("isModalUp", true); })
    .on("hidden.bs.modal", function () { return store.setState("isModalUp", false); });
function init() {
    updateApps();
}
window.addEventListener("keypress", function (e) {
    switch (e.key) {
        case "Enter":
            if (store.getState("isModalUp")) {
                var ev = document.createEvent("Events");
                ev.initEvent("click", true, false);
                modalConfirm.dispatchEvent(ev);
            }
            else if (searchInp == document.activeElement) {
                updateApps(searchInp.value);
            }
            break;
    }
});
function toggleLoader() {
    if (store.getState("loading")) {
        document.querySelector(".loader").classList.remove("d-none");
    }
    else {
        document.querySelector(".loader").classList.add("d-none");
    }
}
function updateModal() {
    console.log(store.getState("isModalUp"));
}
function getOpenExternalButton(hostname, port) {
    var link = location.hostname + ":" + port;
    return "<button class=\"btn btn-secondary\" onclick=\"window.open('" + link + "', '_blank')\"><i class=\"fas fa-external-link-alt fa-2x\"></i><br>Open</button>";
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
    return " <div class=\"card\">\n            <div class=\"card-header\" id=\"heading" + app.id + "\">\n\t\t\t\t<span class=\"float-right " + (running ? "online text-success" : "offline text-danger") + "\">" + (running ? "Online <i class=\"fas fa-globe\"></i>" : "Offline <i class=\"fas fa-times-circle\"></i>") + "</span>\n                <h3 class=\"mb-2\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.name + "\n                </h3>\n                <h6 class=\"mb-0 text-muted\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.repo + "\n\t\t\t\t</h6>\n            </div>\n            <div id=\"collapse" + app.id + "\" class=\"collapse\" aria-labelledby=\"heading" + app.id + "\" data-parent=\"#appContainer\">\n                <div class=\"card-body row p-0\">\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>ID:</span><span>" + app.id + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Name:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Repo:</span><span>" + app.repo + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Root:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Port:</span><span>" + (app.port == 0 ? "none" : app.port) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Hostname:</span><span>" + app.hostname + "</span>\n                        </li>\n                    </ul>\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Deployed:</span><span>" + dateTemplate(app.deployed) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastUpdated:</span><span>" + dateTemplate(app.last_updated) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastRun:</span><span>" + dateTemplate(app.last_run) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Uptime:</span><span>" + (running ? app.uptime.replace(/\.(.*?)s/g, "s") : "offline") + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Runner:</span><span>" + runnerIcon(app.runner) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Pid:</span><span>" + (app.pid == 0 ? "offline" : app.pid) + "</span>\n                        </li>\n                    </ul>\n                </div>\n                <div class=\"card-footer text-right d-flex justify-content-around\">\n                \t" + (running ? getOpenExternalButton("", app.port) : "") + "\n                \t" + (running ? getButton(app.id, "kill") : getButton(app.id, "run")) + "\n    \t\t\t\t" + getButton(app.id, "update") + "\n    \t\t\t\t" + getButton(app.id, "remove") + "\n\t\t\t\t</div>\n            </div>\n        </div>";
}
function updateApps(query) {
    if (query === void 0) { query = ""; }
    var url = baseUrl;
    url.pathname = "/api/find";
    url.search = "?app=" + searchInp.value;
    fetch(url.href).then(function (j) {
        j.json().then(function (res) {
            console.log(res);
            var appsD = res.deployed != null ? res.deployed : [];
            var apps = res.running != null ? res.running : [];
            appContainer.innerHTML = "";
            appsD.forEach(function (a) {
                if (apps.filter(function (app) { return a.id == app.id; }).length == 1) {
                    appContainer.innerHTML += appTemplate(apps.find(function (app) { return a.id == app.id; }), true);
                }
                else {
                    appContainer.innerHTML += appTemplate(a, false);
                }
            });
        }).catch(function (err) { return console.log(err); });
    }).catch(function (err) { return console.log(err); });
}
function doAction(event) {
    var url = baseUrl;
    var btn = event.target;
    var action = btn.attributes.getNamedItem("data-action").value;
    var id = btn.attributes.getNamedItem("data-id").value;
    var data = { app: id };
    url.pathname = "/api/" + action;
    store.setState("loading", true);
    fetch(url.href, {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    })
        .then(function (res) {
        if (res.status == 200) {
            updateApps();
        }
        store.setState("loading", false);
    })
        .catch(function (err) {
        store.setState("loading", false);
        console.log(err);
    });
}
function doForm(event) {
    event.preventDefault();
    var url = baseUrl;
    var repo = document.querySelector("#repo");
    var runner = document.querySelector("#runner");
    var hostname = document.querySelector("#hostname");
    var port = document.querySelector("#port");
    var data = {
        hostname: hostname.value,
        runner: runner.value,
        repo: repo.value.startsWith("https://") ? repo.value : "https://" + repo.value,
        port: port.value,
    };
    url.pathname = "/api/deploy";
    store.setState("loading", true);
    // @ts-ignore
    $("#modalDialog").modal("hide");
    fetch(url.href, {
        method: "POST",
        mode: "cors",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    })
        .then(function (res) {
        if (res.status == 200) {
            // @ts-ignore
            $("#modalDialog").modal("hide");
            updateApps();
        }
        store.setState("loading", false);
    })
        .catch(function (err) {
        console.log(err);
        store.setState("loading", false);
    });
}
init();
