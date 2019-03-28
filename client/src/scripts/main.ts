interface App {
	id: string;
	repo: string;
	name: string;
	root: string;
	port: number;
	hostname: string;
	deployed: string;
	last_updated: string;
	last_run: string;
	uptime: string;
	runner: string;
	pid?: number;
}

type ButtonActions = "run" | "kill" | "update" | "remove";
// interface FindResponse {
// 	running: App[];
// 	deployed: App[];
// }
export type DataStoreTypes = boolean;
export type DataStoreKeys =
	"isModalUp" |
	"isPopUp" |
	"loading";

interface DataStore {
	readonly state: State;
	readonly _state: _State;

	setState(state: DataStoreKeys, value: DataStoreTypes): DataStoreTypes;

	getState(state: DataStoreKeys): DataStoreTypes;

	registerState(state: DataStoreKeys, value: DataStoreTypes): void;

	subscribe(state: DataStoreKeys, actions: Function[]): void;

	getStateObject?(): _State;
}

interface _State {
	[key: string]: _StateProp;
}

interface State {
	[key: string]: DataStoreTypes;
}

interface _StateProp {
	value: DataStoreTypes;
	actions: Function[];
}

function addStyleSheet(rules: string[]) {
	const style = document.createElement("style") as HTMLStyleElement;
	style.appendChild(document.createTextNode(""));
	document.head.append(style);
	for (let i = 0; i < rules.length; i++) {
		(style.sheet as CSSStyleSheet).insertRule(rules[i], i);
	}
}

function initBackdrop(id: string): HTMLElement {
	const bd = document.createElement("div");
	bd.id = id;
	document.body.appendChild(bd);
	return bd;
}

class PopupDialog {
	public confirm: HTMLButtonElement | null;
	public close: HTMLButtonElement | null;
	public popup: HTMLElement | null;
	private readonly backdrop: HTMLElement | null;
	private store: Store;

	constructor(store: Store) {
		this.store = store;
		this.initStates();
		this.initStyleSheet();
		this.backdrop = initBackdrop("popup-backdrop");
		this.popup = null;
		this.confirm = null;
		this.close = null;
	}

	public open(title: string, body: string, cb?: Function) {
		this.createPopup(title, body);
		this.close.addEventListener("click", () => {
			this.destroyPopup();
		});
		if (cb) {
			this.confirm.addEventListener("click", () => {
				cb();
				this.destroyPopup();
			});
			this.confirm.style.display = "inline-block";
		}
		setTimeout(() => {
			this.popup.style.transform = "translateY(10vh)";
		}, 10);
		this.backdrop.style.visibility = "visible";
		this.backdrop.style.opacity = "1";
		this.backdrop.style.top = window.pageYOffset + "px";
		this.store.setState("isPopUp", true);
	}

	public destroyPopup() {
		this.popup.style.transform = "translateY(-10vh)";
		this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
		setTimeout(() => {
			this.confirm.remove();
			this.close.remove();
			this.popup.remove();
			this.popup = null;
			this.confirm = null;
			this.close = null;
			this.backdrop.style.visibility = "hidden";
			this.store.setState("isPopUp", false);
			this.backdrop.style.color = "0";
		}, 100);
	}

	private createPopup(title: string, body: string) {
		const html = `<div id="popup" class="card"><div class="card-header"><h3 class="card-title mb-0">${title}</h3>
						</div><div class="card-body">${body}</div>
						<div class="card-footer">
							<button class="btn btn-danger" id="popupClose"><i class="fas fa-times"></i></button>
							<button class="btn btn-success" id="popupConfirm"><i class="fas fa-check"></i></button>
						</div></div>`;
		this.backdrop.innerHTML += html;
		this.popup = document.querySelector("#popup");
		this.confirm = document.querySelector("#popupConfirm");
		this.close = document.querySelector("#popupClose");
	}

	private initStyleSheet() {
		const rule0 = `#popup-backdrop {
				transition: 100ms all;
				visibility: hidden;
				position: absolute;
				top:0;
				left:0;
				height: 100vh;
				width: 100vw;
				opacity: 1;
				background-color: rgba(0, 0, 0, 0.4);
				z-index: 2000;}`;
		const rule1 = `#popup-backdrop #popup {
				-webkit-transition: 200ms -webkit-transform;
				transition: 200ms -webkit-transform;
				transition: 200ms transform;
				transition: 200ms transform, 200ms -webkit-transform;
				max-width: 600px;
				max-height: 300px;
				margin: 20vh auto;}`;
		const rule2 = `#popup-backdrop #popup .card-body {
			  font-size: 1.5rem;
			  overflow-y: scroll;}`;
		const rule3 = `#popup-backdrop #popup .card-footer {
			  text-align: right;}`;
		const rule4 = `#popup-backdrop #popup #modalConfirm {
			  display: none;}`;
		const rule5 = `#popup-backdrop #popup .card-footer button {
		width: 75px;}`;
		const rules: string[] = [rule0, rule1, rule2, rule3, rule4, rule5];
		addStyleSheet(rules);
	}

	private initStates() {
		this.store.registerState("isPopUp", false);
	}

	public getBackdrop(): HTMLElement {
		return this.backdrop;
	}
}

class Store implements DataStore {
	public readonly _state: _State = {};
	public readonly state: State = {};

	constructor(initialState: State) {
		Object.keys(initialState).forEach(key => {
			this._state[key] = {value: initialState[key], actions: []};
		});
		this.state = initialState;
	}

	public setState(name: DataStoreKeys, state: DataStoreTypes): DataStoreTypes {
		if (Object.keys(this._state).indexOf(name) == -1) {
			throw new Error("State must be registered first");
		} else {
			this.set(name, state);
			if (this._state[name].actions) {
				this._state[name].actions.forEach(action => {
					action();
				});
			}
		}
		return this.state[name];
	}

	public registerState(name: DataStoreKeys, initialState: DataStoreTypes) {
		if (Object.keys(this.state).indexOf(name) != -1) {
			throw new Error("State already exists");
		} else {
			this.set(name, initialState);
		}
	}

	public getState(name: DataStoreKeys): any {
		if (Object.keys(this.state).indexOf(name) == -1) {
			throw new Error("State is not registered - '" + name + "'");
		} else {
			return this.state[name];
		}
	}

	public subscribe(name: DataStoreKeys, actions: Function[]) {
		if (Object.keys(this._state).indexOf(name) == -1) {
			throw new Error("State is not registered");
		} else {
			this._state[name].actions = actions;
		}
	}

	public getStateObject() {
		return this._state;
	}

	private set(name: DataStoreKeys, state: DataStoreTypes) {
		if (Object.keys(this.state).indexOf(name) == -1) {
			this.state[name] = state;
			this._state[name] = {value: state, actions: []};
		} else {
			this._state[name].value = state;
			this.state[name] = state;
		}
	}
}

const initialState: State = {
	isModalUp: false,
	loading: false,
};

const store = new Store(initialState);
const popup = new PopupDialog(store);
store.subscribe("isModalUp", [updateModal]);
store.subscribe("loading", [toggleLoader]);
const baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);
const token = document.cookie.split("; ").filter(e => e.startsWith("Authorization"))[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
let tokenData = jwt_decode(token);
const appContainer = document.querySelector("#appContainer");
const modal = document.querySelector("#modalDialog");
const modalForm = document.querySelector("#modalDialog form") as HTMLFormElement;
const modalConfirm = document.querySelector("#btnModalConfirm");
modalConfirm.addEventListener("click", e => doForm(e));
const modalCancel = document.querySelector("#btnModalCancel");
const searchInp = document.querySelector("#searchInp") as HTMLInputElement;
searchInp.addEventListener("keydown", () => {
	updateApps(searchInp.value);
});
const searchBtn = document.querySelector("#searchBtn");
searchBtn.addEventListener("click", e => {
	updateApps(searchInp.value);
});
const deployBtn = document.querySelector("#deployBtn");
deployBtn.addEventListener("click", e => {
});
$("#modalDialog")
	.on("shown.bs.modal", () => store.setState("isModalUp", true))
	.on("hidden.bs.modal", () => store.setState("isModalUp", false));

function init() {
	updateApps();
}

document.addEventListener("keypress", e => {
	switch (e.key) {
		case "Enter":
			if (store.getState("isModalUp")) {
				const ev = document.createEvent("Events");
				ev.initEvent("click", true, false);
				modalConfirm.dispatchEvent(ev);
			} else if (searchInp == document.activeElement) {
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
	} else {
		document.querySelector(".loader").classList.add("d-none");
	}
}

function updateModal() {
	console.log(store.getState("isModalUp"));
}

function getOpenExternalButton(hostname: string, port: number) {
	let url = window.location.protocol + "//" + hostname;
	if (hostname == "") {
		let sep = window.location.hostname.split(".");
		sep.shift();
		url = window.location.protocol + "//" + sep.join(".") + ":" + port;
	}
	return `<button class="btn btn-secondary" onclick="window.open('${url}', '_blank')"><i class="fas fa-external-link-alt fa-2x"></i><br>Open</button>`;

}

function getButton(id: string, action: ButtonActions): string {
	let icon, color, text;
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

	return `<button class="btn btn-${color}" data-action=\"${action}\" data-id="${id}" onclick=\"doAction(event)\"><i class="fas fa-${icon} fa-2x"></i><br>${text}</button>`;

}

function dateTemplate(dateString: string): string {
	return new Date(dateString).toLocaleString();
}

function runnerIcon(runner: string): string {
	let r = "";
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
	return `<i class="fab fa-${r} fa-2x"></i>`;
}

function appTemplate(app: App, running: boolean): string {
	return ` <div class="card">
            <div class="card-header" id="heading${app.id}">
				<span class="float-right ${running ? "online text-success" : "offline text-danger"}">${running ? "Online <i class=\"fas fa-globe\"></i>" : "Offline <i class=\"fas fa-times-circle\"></i>"}</span>
                <h3 class="mb-2" data-toggle="collapse" data-target="#collapse${app.id}" aria-expanded="false" aria-controls="collapse${app.id}">
					${app.name}
                </h3>
                <h6 class="mb-0 text-muted" data-toggle="collapse" data-target="#collapse${app.id}" aria-expanded="false" aria-controls="collapse${app.id}">
					${app.repo}
				</h6>
            </div>
            <div id="collapse${app.id}" class="collapse" aria-labelledby="heading${app.id}" data-parent="#appContainer">
                <div class="card-body row p-0">
                    <ul class="list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1">
                        <li class="list-group-item d-flex justify-content-between">
                            <span>ID:</span><span>${app.id}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Name:</span><span>${app.name}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Repo:</span><span>${app.repo}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Root:</span><span>${app.name}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Port:</span><span>${app.port == 0 ? "none" : app.port}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Hostname:</span><span>${app.hostname}</span>
                        </li>
                    </ul>
                    <ul class="list-group list-group-flush col-lg-6 col-md-12 pl-1 pr-1">
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Deployed:</span><span>${dateTemplate(app.deployed)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>LastUpdated:</span><span>${dateTemplate(app.last_updated)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>LastRun:</span><span>${dateTemplate(app.last_run)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Uptime:</span><span>${running ? app.uptime.replace(/\.(.*?)s/g, "s") : "offline"}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Runner:</span><span>${runnerIcon(app.runner)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Pid:</span><span>${app.pid == 0 ? "offline" : app.pid}</span>
                        </li>
                    </ul>
                </div>
                <div class="card-footer text-right d-flex justify-content-around">
                	${running ? getOpenExternalButton(app.hostname, app.port) : ""}
                	${running ? getButton(app.id, "kill") : getButton(app.id, "run")}
    				${getButton(app.id, "update")}
    				${getButton(app.id, "remove")}
				</div>
            </div>
        </div>`;
}

function updateApps(query: string = "") {
	const url = baseUrl;
	url.pathname = "/api/find";
	url.search = "?app=" + searchInp.value;
	fetch(url.href).then(j => {
		j.json().then(res => {
			console.log(res);
			const appsD: App[] = res.deployed != null ? res.deployed : [];
			const apps: App[] = res.running != null ? res.running : [];
			appContainer.innerHTML = "";
			appsD.forEach(a => {
				if (apps.filter(app => a.id == app.id).length == 1) {
					appContainer.innerHTML += appTemplate(apps.find(app => a.id == app.id), true);
				} else {
					appContainer.innerHTML += appTemplate(a, false);
				}
			});
		}).catch(err => console.log(err));
	}).catch(err => console.log(err));
}

function doAction(event: Event) {
	popup.open("Warning", "Are you sure?", () => {
		const url = baseUrl;
		const btn = event.target as HTMLButtonElement;
		const action = btn.attributes.getNamedItem("data-action").value;
		const id = btn.attributes.getNamedItem("data-id").value;
		const data = {app: id};
		url.pathname = `/api/${action}`;
		store.setState("loading", true);
		fetch(url.href, {
			method: "POST", // *GET, POST, PUT, DELETE, etc.
			mode: "cors", // no-cors, cors, *same-origin
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data), // body data type must match "Content-Type" header
		})
			.then(res => {
				if (res.status == 200) {
					updateApps();
				}
				store.setState("loading", false);
			})
			.catch(err => {
				store.setState("loading", false);
				console.log(err);
			});
	});
}

function doForm(event: Event) {
	popup.open("Warning", "Are you sure?", () => {
		event.preventDefault();
		const url = baseUrl;
		const repo = document.querySelector("#repo") as HTMLInputElement;
		const runner = document.querySelector("#runner") as HTMLInputElement;
		const hostname = document.querySelector("#hostname") as HTMLInputElement;
		const port = document.querySelector("#port") as HTMLInputElement;
		const data = {
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
			method: "POST", // *GET, POST, PUT, DELETE, etc.
			mode: "cors", // no-cors, cors, *same-origin
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data), // body data type must match "Content-Type" header
		})
			.then(res => {
				if (res.status == 200) {
					// @ts-ignore
					$("#modalDialog").modal("hide");
					updateApps();
				}
				store.setState("loading", false);
			})
			.catch(err => {
				console.log(err);
				store.setState("loading", false);
			});
	});
}

init();
