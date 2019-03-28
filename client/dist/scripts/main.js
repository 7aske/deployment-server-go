"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
function addStyleSheet(rules) {
    var style = document.createElement("style");
    style.appendChild(document.createTextNode(""));
    document.head.append(style);
    for (var i = 0; i < rules.length; i++) {
        style.sheet.insertRule(rules[i], i);
    }
}
function initBackdrop(id) {
    var bd = document.createElement("div");
    bd.id = id;
    document.body.appendChild(bd);
    return bd;
}
var PopupDialog = /** @class */ (function () {
    function PopupDialog(store) {
        this.store = store;
        this.initStates();
        this.initStyleSheet();
        this.backdrop = initBackdrop("popup-backdrop");
        this.popup = null;
        this.confirm = null;
        this.close = null;
    }
    PopupDialog.prototype.open = function (title, body, cb) {
        var _this = this;
        this.createPopup(title, body);
        this.close.addEventListener("click", function () {
            _this.destroyPopup();
        });
        if (cb) {
            this.confirm.addEventListener("click", function () {
                cb();
                _this.destroyPopup();
            });
            this.confirm.style.display = "inline-block";
        }
        setTimeout(function () {
            _this.popup.style.transform = "translateY(10vh)";
        }, 10);
        this.backdrop.style.visibility = "visible";
        this.backdrop.style.opacity = "1";
        this.backdrop.style.top = window.pageYOffset + "px";
        this.store.setState("isPopUp", true);
    };
    PopupDialog.prototype.destroyPopup = function () {
        var _this = this;
        this.popup.style.transform = "translateY(-10vh)";
        this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
        setTimeout(function () {
            _this.confirm.remove();
            _this.close.remove();
            _this.popup.remove();
            _this.popup = null;
            _this.confirm = null;
            _this.close = null;
            _this.backdrop.style.visibility = "hidden";
            _this.store.setState("isPopUp", false);
            _this.backdrop.style.color = "0";
        }, 100);
    };
    PopupDialog.prototype.createPopup = function (title, body) {
        var html = "<div id=\"popup\" class=\"card\"><div class=\"card-header\"><h3 class=\"card-title mb-0\">" + title + "</h3>\n\t\t\t\t\t\t</div><div class=\"card-body\">" + body + "</div>\n\t\t\t\t\t\t<div class=\"card-footer\">\n\t\t\t\t\t\t\t<button class=\"btn btn-danger\" id=\"popupClose\"><i class=\"fas fa-times\"></i></button>\n\t\t\t\t\t\t\t<button class=\"btn btn-success\" id=\"popupConfirm\"><i class=\"fas fa-check\"></i></button>\n\t\t\t\t\t\t</div></div>";
        this.backdrop.innerHTML += html;
        this.popup = document.querySelector("#popup");
        this.confirm = document.querySelector("#popupConfirm");
        this.close = document.querySelector("#popupClose");
    };
    PopupDialog.prototype.initStyleSheet = function () {
        var rule0 = "#popup-backdrop {\n\t\t\t\ttransition: 100ms all;\n\t\t\t\tvisibility: hidden;\n\t\t\t\tposition: absolute;\n\t\t\t\ttop:0;\n\t\t\t\tleft:0;\n\t\t\t\theight: 100vh;\n\t\t\t\twidth: 100vw;\n\t\t\t\topacity: 1;\n\t\t\t\tbackground-color: rgba(0, 0, 0, 0.4);\n\t\t\t\tz-index: 2000;}";
        var rule1 = "#popup-backdrop #popup {\n\t\t\t\t-webkit-transition: 200ms -webkit-transform;\n\t\t\t\ttransition: 200ms -webkit-transform;\n\t\t\t\ttransition: 200ms transform;\n\t\t\t\ttransition: 200ms transform, 200ms -webkit-transform;\n\t\t\t\tmax-width: 600px;\n\t\t\t\tmax-height: 300px;\n\t\t\t\tmargin: 20vh auto;}";
        var rule2 = "#popup-backdrop #popup .card-body {\n\t\t\t  font-size: 1.5rem;\n\t\t\t  overflow-y: scroll;}";
        var rule3 = "#popup-backdrop #popup .card-footer {\n\t\t\t  text-align: right;}";
        var rule4 = "#popup-backdrop #popup #modalConfirm {\n\t\t\t  display: none;}";
        var rule5 = "#popup-backdrop #popup .card-footer button {\n\t\twidth: 75px;}";
        var rules = [rule0, rule1, rule2, rule3, rule4, rule5];
        addStyleSheet(rules);
    };
    PopupDialog.prototype.initStates = function () {
        this.store.registerState("isPopUp", false);
    };
    PopupDialog.prototype.getBackdrop = function () {
        return this.backdrop;
    };
    return PopupDialog;
}());
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
var popup = new PopupDialog(store);
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
document.addEventListener("keypress", function (e) {
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
        case "Escape":
            if (store.getState("isPopUp")) {
                store.setState("isPopUp", false);
                popup.destroyPopup();
            }
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
    var url = window.location.protocol + "//" + hostname;
    if (hostname == "") {
        var sep = window.location.hostname.split(".");
        sep.shift();
        url = window.location.protocol + "//" + sep.join(".") + ":" + port;
    }
    return "<button class=\"btn btn-secondary\" onclick=\"window.open('" + url + "', '_blank')\"><i class=\"fas fa-external-link-alt fa-2x\"></i><br>Open</button>";
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
        case "python":
            r = "python";
            break;
    }
    return "<i class=\"fab fa-" + r + " fa-2x\"></i>";
}
function appTemplate(app, running) {
    return " <div class=\"card\">\n            <div class=\"card-header\" id=\"heading" + app.id + "\">\n\t\t\t\t<span class=\"float-right " + (running ? "online text-success" : "offline text-danger") + "\">" + (running ? "Online <i class=\"fas fa-globe\"></i>" : "Offline <i class=\"fas fa-times-circle\"></i>") + "</span>\n                <h3 class=\"mb-2\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.name + "\n                </h3>\n                <h6 class=\"mb-0 text-muted\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.repo + "\n\t\t\t\t</h6>\n            </div>\n            <div id=\"collapse" + app.id + "\" class=\"collapse\" aria-labelledby=\"heading" + app.id + "\" data-parent=\"#appContainer\">\n                <div class=\"card-body row p-0\">\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>ID:</span><span>" + app.id + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Name:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Repo:</span><span>" + app.repo + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Root:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Port:</span><span>" + (app.port == 0 ? "none" : app.port) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Hostname:</span><span>" + app.hostname + "</span>\n                        </li>\n                    </ul>\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Deployed:</span><span>" + dateTemplate(app.deployed) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastUpdated:</span><span>" + dateTemplate(app.last_updated) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastRun:</span><span>" + dateTemplate(app.last_run) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Uptime:</span><span>" + (running ? app.uptime.replace(/\.(.*?)s/g, "s") : "offline") + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Runner:</span><span>" + runnerIcon(app.runner) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Pid:</span><span>" + (app.pid == 0 ? "offline" : app.pid) + "</span>\n                        </li>\n                    </ul>\n                </div>\n                <div class=\"card-footer text-right d-flex justify-content-around\">\n                \t" + (running ? getOpenExternalButton(app.hostname, app.port) : "") + "\n                \t" + (running ? getButton(app.id, "kill") : getButton(app.id, "run")) + "\n    \t\t\t\t" + getButton(app.id, "update") + "\n    \t\t\t\t" + getButton(app.id, "remove") + "\n\t\t\t\t</div>\n            </div>\n        </div>";
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
    popup.open("Warning", "Are you sure?", function () {
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
    });
}
function doForm(event) {
    popup.open("Warning", "Are you sure?", function () {
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
    });
}
init();
