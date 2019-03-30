"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : new P(function (resolve) { resolve(result.value); }).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
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
        this.confirmBtn = null;
        this.closeBtn = null;
    }
    PopupDialog.prototype.confirm = function () {
        var ev = document.createEvent("Events");
        ev.initEvent("click", true, false);
        this.confirmBtn.dispatchEvent(ev);
    };
    PopupDialog.prototype.cancel = function () {
        var ev = document.createEvent("Events");
        ev.initEvent("click", true, false);
        this.closeBtn.dispatchEvent(ev);
        this.store.setState("isPopUp", false);
    };
    PopupDialog.prototype.open = function (title, body, cb) {
        var _this = this;
        this.createPopup(title, body);
        this.closeBtn.addEventListener("click", function () {
            _this.destroyPopup();
        });
        if (cb) {
            this.confirmBtn.addEventListener("click", function () {
                cb();
                _this.destroyPopup();
            });
            this.confirmBtn.style.display = "inline-block";
        }
        setTimeout(function () {
            _this.popup.style.transform = "translateY(10vh)";
        }, 10);
        this.backdrop.style.visibility = "visible";
        this.backdrop.style.opacity = "1";
        this.backdrop.style.height = document.body.offsetHeight + "px";
        this.popup.style.top = window.pageYOffset + "px";
        this.store.setState("isPopUp", true);
    };
    PopupDialog.prototype.destroyPopup = function () {
        var _this = this;
        this.popup.style.transform = "translateY(-10vh)";
        this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
        setTimeout(function () {
            _this.confirmBtn.remove();
            _this.closeBtn.remove();
            _this.popup.remove();
            _this.popup = null;
            _this.confirmBtn = null;
            _this.closeBtn = null;
            _this.backdrop.style.visibility = "hidden";
            _this.store.setState("isPopUp", false);
            _this.backdrop.style.color = "0";
        }, 100);
    };
    PopupDialog.prototype.createPopup = function (title, body) {
        var html = "<div id=\"popup\" class=\"card\"><div class=\"card-header\"><h3 class=\"card-title mb-0\">" + title + "</h3>\n\t\t\t\t\t\t</div><div class=\"card-body\">" + body + "</div>\n\t\t\t\t\t\t<div class=\"card-footer\">\n\t\t\t\t\t\t\t<button class=\"btn btn-danger\" id=\"popupClose\"><i class=\"fas fa-times\"></i></button>\n\t\t\t\t\t\t\t<button class=\"btn btn-success\" id=\"popupConfirm\"><i class=\"fas fa-check\"></i></button>\n\t\t\t\t\t\t</div></div>";
        this.backdrop.innerHTML += html;
        this.popup = document.querySelector("#popup");
        this.confirmBtn = document.querySelector("#popupConfirm");
        this.closeBtn = document.querySelector("#popupClose");
    };
    PopupDialog.prototype.initStyleSheet = function () {
        var rule0 = "#popup-backdrop {\n\t\t\t\ttransition: 100ms all;\n\t\t\t\tvisibility: hidden;\n\t\t\t\tposition: absolute;\n\t\t\t\ttop:0;\n\t\t\t\tleft:0;\n\t\t\t\theight: 100vh;\n\t\t\t\twidth: 100vw;\n\t\t\t\topacity: 1;\n\t\t\t\tbackground-color: rgba(0, 0, 0, 0.4);\n\t\t\t\tz-index: 2000;}";
        var rule1 = "#popup-backdrop #popup {\n\t\t\t\t-webkit-transition: 200ms -webkit-transform;\n\t\t\t\ttransition: 200ms -webkit-transform;\n\t\t\t\ttransition: 200ms transform;\n\t\t\t\ttransition: 200ms transform, 200ms -webkit-transform;\n\t\t\t\tmax-width: 600px;\n\t\t\t\tmax-height: 300px;\n\t\t\t\tmargin: 20vh auto;}";
        var rule2 = "#popup-backdrop #popup .card-body {\n\t\t\t\tfont-size: 1.5rem;\n\t\t\t\toverflow-y: auto;}";
        var rule3 = "#popup-backdrop #popup .card-footer {\n\t\t\t\ttext-align: right;}";
        var rule4 = "#popup-backdrop #popup #popupConfirm {\n\t\t\t\tdisplay: none;}";
        var rule5 = "#popup-backdrop #popup .card-footer button {\n\t\t\t\twidth: 75px;}";
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
var Modal = /** @class */ (function () {
    function Modal(store) {
        this.store = store;
        this.initStyleSheets();
        this.initStates();
        this.modal = document.createElement("section");
        this.backdrop = initBackdrop("modal-backdrop");
        this.closeBtn = null;
    }
    Modal.prototype.createModal = function (header, body) {
        this.backdrop.innerHTML = "<div id=\"modal\" class=\"card\"><div class=\"card-header\"><h5 class=\"card-title mb-0\">" + (header ? header : "") + "</h5>\n\t\t\t\t\t\t</div><div class=\"card-body\">" + (body ? body : "") + "</div>\n\t\t\t\t\t\t<div class=\"card-footer pl-3\">\n\t\t\t\t\t\t\t<button class=\"btn btn-secondary\" id=\"modalClose\">Close</button>\n\t\t\t\t\t\t</div></div>";
        this.closeBtn = document.querySelector("#modalClose");
        this.closeBtn = document.querySelector("#modalClose");
        this.modal = document.querySelector("#modal");
    };
    Modal.prototype.open = function (header, body, cb) {
        var _this = this;
        this.createModal(header, body);
        this.closeBtn.addEventListener("click", function () { return _this.destroyModal(); });
        // this.confirmBtn.addEventListener("click", () => cb());
        setTimeout(function () {
            _this.modal.style.transform = "translateY(10vh)";
        }, 10);
        this.backdrop.style.visibility = "visible";
        this.backdrop.style.opacity = "1";
        this.store.setState("isModalUp", true);
        this.up = true;
        // this.backdrop.style.top = window.pageYOffset + "px";
        this.backdrop.style.height = document.body.offsetHeight + "px";
        this.modal.style.top = window.pageYOffset + "px";
    };
    Modal.prototype.close = function () {
        this.destroyModal();
        this.up = false;
        store.setState("currentApp", null);
    };
    Modal.prototype.destroyModal = function () {
        var _this = this;
        this.modal.style.transform = "translateY(-10vh)";
        this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
        setTimeout(function () {
            if (_this.closeBtn)
                _this.closeBtn.remove();
            if (_this.modal)
                _this.modal.remove();
            if (_this.script)
                _this.script.remove();
            _this.modal = null;
            _this.closeBtn = null;
            _this.script = null;
            _this.backdrop.style.visibility = "hidden";
            _this.store.setState("isModalUp", false);
            _this.backdrop.style.color = "0";
        }, 100);
    };
    Modal.prototype.runScripts = function (src) {
        this.script = document.createElement("script");
        this.script.src = src;
        this.backdrop.appendChild(this.script);
    };
    Modal.prototype.initStates = function () {
        if (!this.store.hasState("isModalUp"))
            this.store.registerState("isModalUp", false);
    };
    Modal.prototype.initStyleSheets = function () {
        var rule0 = "#modal-backdrop {\n\t\t\ttransition: 100ms all;\n\t\t\tvisibility: hidden;\n\t\t\tposition: absolute;\n\t\t\ttop: 0;\n\t\t\tleft:0;\n\t\t\theight: 100vh;\n\t\t\twidth: 100vw;\n\t\t\topacity: 1;\n\t\t\tbackground-color: rgba(0, 0, 0, 0.4);\n\t\t\tz-index: 1500;\n\t\t\tpadding: 20px;\n\n\t\t}";
        var rule1 = "#modal-backdrop #modal {\n\t\t\t-webkit-transition: 200ms -webkit-transform;\n\t\t\ttransition: 200ms -webkit-transform;\n\t\t\ttransition: 200ms transform;\n\t\t\ttransition: 200ms transform, 200ms -webkit-transform;\n\t\t\tmax-width: 800px;\n\t\t\tmin-height: 400px;\n\t\t\tmargin: auto;\t\t\t\n\t\t}";
        addStyleSheet([rule0, rule1]);
    };
    Modal.prototype.getBackdrop = function () {
        return this.backdrop;
    };
    Modal.prototype.getModal = function () {
        return this.modal;
    };
    return Modal;
}());
exports.Modal = Modal;
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
    Store.prototype.hasState = function (state) {
        if (Object.keys(this.state).indexOf(state) == -1) {
            return false;
        }
        else {
            return true;
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
    runningApps: [],
    deployedApps: [],
    currentApp: null,
};
var store = new Store(initialState);
var popup = new PopupDialog(store);
store.subscribe("isModalUp", [updateModal]);
store.subscribe("loading", [toggleLoader]);
var baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);
var token = document.cookie.split("; ").filter(function (e) { return e.startsWith("Authorization"); })[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
var tokenData = jwt_decode(token);
var modal = new Modal(store);
var appContainer = document.querySelector("#appContainer");
var deployDialog = document.querySelector("#deployDialog");
var deployDialogForm = document.querySelector("#deployDialog form");
var deployDialogConfirm = document.querySelector("#btnModalConfirm");
deployDialogConfirm.addEventListener("click", function (e) { return doForm(e); });
var deployDialogCancel = document.querySelector("#btnModalCancel");
var searchInp = document.querySelector("#searchInp");
searchInp.addEventListener("keydown", function (e) {
    if (e.key == "Backspace" && searchInp.value.length == 1) {
        updateApps();
    }
    updateApps(searchInp.value);
});
var searchBtn = document.querySelector("#searchBtn");
searchBtn.addEventListener("click", function (e) {
    updateApps(searchInp.value);
});
var deployBtn = document.querySelector("#deployBtn");
deployBtn.addEventListener("click", function (e) {
});
$("#deployDialog")
    .on("shown.bs.modal", function () { return store.setState("isModalUp", true); })
    .on("hidden.bs.modal", function () { return store.setState("isModalUp", false); });
function init() {
    updateApps();
}
document.addEventListener("keydown", function (e) {
    switch (e.key) {
        case "Enter":
            if (store.getState("isPopUp")) {
                popup.confirm();
            }
            else if (store.getState("isModalUp")) {
                deployDialogConfirm.click();
            }
            else if (searchInp == document.activeElement) {
                updateApps(searchInp.value);
            }
            break;
        case "Escape":
            if (store.getState("isPopUp")) {
                popup.cancel();
            }
            else if (store.getState("isModalUp")) {
                if (modal.up) {
                    modal.close();
                }
                else {
                    deployDialogCancel.click();
                }
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
        // let sep = window.location.hostname.split(".");
        // sep.shift();
        // url = window.location.protocol + "//" + sep.join(".") + ":" + port;
        url += window.location.hostname + ":" + port;
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
        case "settings":
            color = "secondary";
            icon = "cogs";
            text = "Settings";
            break;
    }
    return "<button class=\"btn btn-" + color + "\" data-action=\"" + action + "\" data-id=\"" + id + "\" onclick=\"doAction(event)\"><i class=\"fas fa-" + icon + " fa-2x\"></i><br>" + text + "</button>";
}
function dateTemplate(dateString) {
    return new Date(dateString).toLocaleString();
}
function getRunnerIcon(runner) {
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
    return " <div class=\"card\">\n            <div class=\"card-header\" id=\"heading" + app.id + "\">\n\t\t\t\t<span class=\"float-right " + (running ? "online text-success" : "offline text-danger") + "\">" + (running ? "Online <i class=\"fas fa-globe\"></i>" : "Offline <i class=\"fas fa-times-circle\"></i>") + "</span>\n                <h3 class=\"mb-2\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.name + "\n                </h3>\n                <h6 class=\"mb-0 text-muted\" data-toggle=\"collapse\" data-target=\"#collapse" + app.id + "\" aria-expanded=\"false\" aria-controls=\"collapse" + app.id + "\">\n\t\t\t\t\t" + app.repo + "\n\t\t\t\t</h6>\n            </div>\n            <div id=\"collapse" + app.id + "\" class=\"collapse\" aria-labelledby=\"heading" + app.id + "\" data-parent=\"#appContainer\">\n                <div class=\"card-body row p-0\">\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>ID:</span><span>" + app.id + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Name:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Repo:</span><span>" + app.repo + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Root:</span><span>" + app.name + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Port:</span><span>" + (app.port == 0 ? "none" : app.port) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Hostname:</span><span>" + (app.hostname == "" ? window.location.hostname + ":" + app.port : app.hostname) + "</span>\n                        </li>\n                    </ul>\n                    <ul class=\"list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1\">\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Deployed:</span><span>" + dateTemplate(app.deployed) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastUpdated:</span><span>" + dateTemplate(app.last_updated) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>LastRun:</span><span>" + dateTemplate(app.last_run) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Uptime:</span><span>" + (running ? app.uptime.replace(/\.(.*?)s/g, "s") : "offline") + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Runner:</span><span>" + getRunnerIcon(app.runner) + "</span>\n                        </li>\n                        <li class=\"list-group-item d-flex justify-content-between\">\n                            <span>Pid:</span><span>" + (app.pid == -1 ? "offline" : app.pid) + "</span>\n                        </li>\n                    </ul>\n                </div>\n                <div class=\"card-footer text-right d-flex justify-content-around\">\n                \t" + (running ? getOpenExternalButton(app.hostname, app.port) : getButton(app.id, "settings")) + "\n                \t" + (running ? getButton(app.id, "kill") : getButton(app.id, "run")) + "\n    \t\t\t\t" + getButton(app.id, "update") + "\n    \t\t\t\t" + getButton(app.id, "remove") + "\n\t\t\t\t</div>\n            </div>\n        </div>";
}
function settingsTemplate(app) {
    return "<div class=\"row\"><div class=\"col\"></div><form class=\"col-md-8\">\n\t\t\t<div class=\"input-group input-group-sm mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<span class=\"input-group-text\">ID</span>\n\t\t\t\t</div>\n\t\t\t\t<input readonly value=\"" + app.id + "\" type=\"text\" class=\"form-control\" aria-label=\"Small\">\n\t\t\t</div>\n\t\t\t<div class=\"input-group input-group-sm mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<span class=\"input-group-text\">Name</span>\n\t\t\t\t</div>\n\t\t\t\t<input readonly value=\"" + app.name + "\" type=\"text\" class=\"form-control\" aria-label=\"Small\">\n\t\t\t</div>\n\t\t\t<div class=\"input-group input-group-sm mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<span class=\"input-group-text\">Repository</span>\n\t\t\t\t</div>\n\t\t\t\t<input readonly value=\"" + app.repo + "\" type=\"text\" class=\"form-control\" aria-label=\"Small\">\n\t\t\t</div>\n\t\t\t<div class=\"input-group mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<label class=\"input-group-text\" for=\"runner\">Runner</label>\n\t\t\t\t</div>\n\t\t\t\t<select class=\"custom-select\" name=\"runner\" id=\"runnerSettings\">\n\t\t\t\t\t<option " + (app.runner == "node" ? "selected" : "") + " value=\"node\">Node</option>\n\t\t\t\t\t<option " + (app.runner == "web" ? "selected" : "") + " value=\"web\">Web</option>\n\t\t\t\t\t<option " + (app.runner == "python" ? "selected" : "") + " value=\"python\">Python</option>\n\t\t\t\t</select>\n\t\t\t</div>\n\t\t\t<div class=\"input-group input-group-sm mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<span class=\"input-group-text\">Hostname</span>\n\t\t\t\t</div>\n\t\t\t\t<input value=\"" + app.hostname + "\" type=\"text\" name=\"hostname\" id=\"hostnameSettings\" class=\"form-control\" aria-label=\"Small\">\n\t\t\t</div>\n\t\t\t<div class=\"input-group input-group-sm mb-3\">\n\t\t\t\t<div class=\"input-group-prepend\">\n\t\t\t\t\t<span class=\"input-group-text\">Port</span>\n\t\t\t\t</div>\n\t\t\t\t<input value=\"" + app.port + "\" type=\"text\" name=\"port\" id=\"portSettings\" class=\"form-control\" aria-label=\"Small\">\n\t\t\t</div>\n\t\t</form><div class=\"col\"></div></div>\n\t\t<button class=\"btn btn-success\" type=\"button\" data-id=\"" + app.id + "\" data-action=\"settings\" onclick=\"doModalForm(event)\"\">Update</button>";
}
function updateApps(query) {
    if (query === void 0) { query = ""; }
    var url = baseUrl;
    url.pathname = "/api/find";
    url.search = "?app=" + query;
    fetch(url.href).then(function (j) {
        j.json().then(function (res) {
            var appsD = res.deployed != null ? res.deployed : [];
            var apps = res.running != null ? res.running : [];
            if (query == "") {
                store.setState("deployedApps", appsD);
                store.setState("runningApps", apps);
                console.log(store);
            }
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
    if (action != "settings") {
        popup.open("Warning", "Are you sure?", function () {
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
    else {
        var apps = store.getState("deployedApps");
        var app = apps.find(function (a) { return a.id == id; });
        store.setState("currentApp", app);
        modal.open(app.name, settingsTemplate(app));
    }
}
function doModalForm(event) {
    return __awaiter(this, void 0, void 0, function () {
        var _this = this;
        return __generator(this, function (_a) {
            popup.open("Warning", "Are you sure?", function () {
                event.preventDefault();
                var btn = event.target;
                var action = btn.attributes.getNamedItem("data-action").value;
                var id = btn.attributes.getNamedItem("data-id").value;
                var app = store.getState("currentApp");
                var url = baseUrl;
                var edits = {
                    "runner": document.querySelector("#runnerSettings").value,
                    "hostname": document.querySelector("#hostnameSettings").value,
                    "port": document.querySelector("#portSettings").value,
                };
                var settings = {};
                Object.keys(edits).forEach(function (key) {
                    if (edits[key] != String(app[key])) {
                        settings[key] = edits[key];
                    }
                });
                var data = { id: id, settings: settings };
                url.pathname = "/api/" + action;
                console.log(action);
                store.setState("loading", true);
                modal.close();
                fetch(url.href, {
                    method: "POST",
                    mode: "cors",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(data),
                })
                    .then(function (res) { return __awaiter(_this, void 0, void 0, function () {
                    var response_1;
                    return __generator(this, function (_a) {
                        switch (_a.label) {
                            case 0:
                                console.log(res);
                                if (!(res.status == 200)) return [3 /*break*/, 1];
                                updateApps();
                                return [3 /*break*/, 3];
                            case 1:
                                if (!(res.status == 500)) return [3 /*break*/, 3];
                                return [4 /*yield*/, res.json()];
                            case 2:
                                response_1 = _a.sent();
                                setTimeout(function () {
                                    popup.open("Error", response_1.message);
                                }, 200);
                                console.log(response_1);
                                _a.label = 3;
                            case 3:
                                store.setState("loading", false);
                                return [2 /*return*/];
                        }
                    });
                }); })
                    .catch(function (err) {
                    console.log(err);
                    store.setState("loading", false);
                });
            });
            return [2 /*return*/];
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
        $("#deployDialog").modal("hide");
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
                $("#deployDialog").modal("hide");
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
