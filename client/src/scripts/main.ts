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

	[key: string]: string | number
}

interface Edits {
	runner?: string;
	port?: string;
	hostname?: string;

	[key: string]: string;
}

type ButtonActions = "run" | "kill" | "update" | "remove" | "settings";
// interface FindResponse {
// 	running: App[];
// 	deployed: App[];
// }
export type DataStoreTypes = boolean | App | App[];
export type DataStoreKeys =
	"isModalUp" |
	"isPopUp" |
	"runningApps" |
	"deployedApps" |
	"currentApp" |
	"loading";

interface DataStore {
	readonly state: State;
	readonly _state: _State;

	setState(state: DataStoreKeys, value: DataStoreTypes): DataStoreTypes;

	getState(state: DataStoreKeys): DataStoreTypes;

	hasState(state: string): boolean

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
	private confirmBtn: HTMLButtonElement | null;
	private closeBtn: HTMLButtonElement | null;
	private popup: HTMLElement | null;
	private readonly backdrop: HTMLElement | null;
	private store: Store;

	constructor(store: Store) {
		this.store = store;
		this.initStates();
		PopupDialog.initStyleSheet();
		this.backdrop = initBackdrop("popup-backdrop");
		this.popup = null;
		this.confirmBtn = null;
		this.closeBtn = null;
	}

	public confirm() {
		const ev = document.createEvent("Events");
		ev.initEvent("click", true, false);
		this.confirmBtn.dispatchEvent(ev);
	}

	public cancel() {
		const ev = document.createEvent("Events");
		ev.initEvent("click", true, false);
		this.closeBtn.dispatchEvent(ev);
		this.store.setState("isPopUp", false);
	}

	public open(title: string, body: string, cb?: Function) {
		this.createPopup(title, body);
		this.closeBtn.addEventListener("click", () => {
			this.destroyPopup();
		});
		if (cb) {
			this.confirmBtn.addEventListener("click", () => {
				cb();
				this.destroyPopup();
			});
			this.confirmBtn.style.display = "inline-block";
		}
		setTimeout(() => {
			this.popup.style.transform = "translateY(10vh)";
		}, 10);
		this.backdrop.style.visibility = "visible";
		this.backdrop.style.opacity = "1";
		this.backdrop.style.height = document.body.offsetHeight + "px";
		this.popup.style.top = window.pageYOffset + "px";
		this.store.setState("isPopUp", true);
	}

	public destroyPopup() {
		this.popup.style.transform = "translateY(-10vh)";
		this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
		setTimeout(() => {
			this.confirmBtn.remove();
			this.closeBtn.remove();
			this.popup.remove();
			this.popup = null;
			this.confirmBtn = null;
			this.closeBtn = null;
			this.backdrop.style.visibility = "hidden";
			this.store.setState("isPopUp", false);
			this.backdrop.style.color = "0";
		}, 100);
	}

	private createPopup(title: string, body: string) {
		this.backdrop.innerHTML = `<div id="popup" class="card"><div class="card-header"><h3 class="card-title mb-0">${title}</h3>
						</div><div class="card-body">${body}</div>
						<div class="card-footer">
							<button class="btn btn-danger" id="popupClose"><i class="fas fa-times"></i></button>
							<button class="btn btn-success" id="popupConfirm"><i class="fas fa-check"></i></button>
						</div></div>`;
		this.popup = document.querySelector("#popup");
		this.confirmBtn = document.querySelector("#popupConfirm");
		this.closeBtn = document.querySelector("#popupClose");
	}

	private static initStyleSheet() {
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
				overflow-y: auto;}`;
		const rule3 = `#popup-backdrop #popup .card-footer {
				text-align: right;}`;
		const rule4 = `#popup-backdrop #popup #popupConfirm {
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

export class Modal {
	private readonly backdrop: HTMLElement;
	private readonly store: Store;
	private modal: HTMLElement;
	private closeBtn: HTMLButtonElement;
	private script: HTMLScriptElement;
	public up: boolean;

	constructor(store: Store) {
		this.store = store;
		Modal.initStyleSheets();
		this.initStates();
		this.modal = document.createElement("section");
		this.backdrop = initBackdrop("modal-backdrop");
		this.closeBtn = null;
	}

	private createModal(header?: string, body?: string) {
		this.backdrop.innerHTML = `<div id="modal" class="card"><div class="card-header"><h5 class="card-title mb-0">${header ? header : ""}</h5>
						</div><div class="card-body">${body ? body : ""}</div>
						<div class="card-footer pl-3">
							<button class="btn btn-secondary" id="modalClose">Close</button>
						</div></div>`;
		this.closeBtn = document.querySelector("#modalClose") as HTMLButtonElement;
		this.closeBtn = document.querySelector("#modalClose") as HTMLButtonElement;
		this.modal = document.querySelector("#modal");
	}

	public open(header?: string, body?: string, cb?: Function) {
		this.createModal(header, body);
		this.closeBtn.addEventListener("click", () => this.destroyModal());
		// this.confirmBtn.addEventListener("click", () => cb());
		setTimeout(() => {
			this.modal.style.transform = "translateY(10vh)";
		}, 10);
		this.backdrop.style.visibility = "visible";
		this.backdrop.style.opacity = "1";
		this.store.setState("isModalUp", true);
		this.up = true;
		// this.backdrop.style.top = window.pageYOffset + "px";
		this.backdrop.style.height = document.body.offsetHeight + "px";
		this.modal.style.top = window.pageYOffset + "px";
	}

	public close() {
		this.destroyModal();
		this.up = false;
		store.setState("currentApp", null);
	}

	private destroyModal() {
		this.modal.style.transform = "translateY(-10vh)";
		this.backdrop.style.backgroundColor = "background-color: rgba(0, 0, 0, 0)";
		setTimeout(() => {
			if (this.closeBtn)
				this.closeBtn.remove();
			if (this.modal)
				this.modal.remove();
			if (this.script)
				this.script.remove();
			this.modal = null;
			this.closeBtn = null;
			this.script = null;
			this.backdrop.style.visibility = "hidden";
			this.store.setState("isModalUp", false);
			this.backdrop.style.color = "0";
		}, 100);
	}

	public runScripts(src: string) {
		this.script = document.createElement("script");
		this.script.src = src;
		this.backdrop.appendChild(this.script);
	}

	private initStates() {
		if (!this.store.hasState("isModalUp"))
			this.store.registerState("isModalUp", false);
	}

	private static initStyleSheets() {
		const rule0 = `#modal-backdrop {
			transition: 100ms all;
			visibility: hidden;
			position: absolute;
			top: 0;
			left:0;
			height: 100vh;
			width: 100vw;
			opacity: 1;
			background-color: rgba(0, 0, 0, 0.4);
			z-index: 1500;
			padding: 20px;

		}`;
		const rule1 = `#modal-backdrop #modal {
			-webkit-transition: 200ms -webkit-transform;
			transition: 200ms -webkit-transform;
			transition: 200ms transform;
			transition: 200ms transform, 200ms -webkit-transform;
			max-width: 800px;
			min-height: 400px;
			margin: auto;			
		}`;
		addStyleSheet([rule0, rule1]);
	}

	public getBackdrop(): HTMLElement {
		return this.backdrop;
	}

	public getModal(): HTMLElement {
		return this.modal;
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

	public hasState(state: string): boolean {
		return Object.keys(this.state).indexOf(state) != -1;
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
	runningApps: [],
	deployedApps: [],
	currentApp: null,
};

const token = document.cookie.split("; ").filter(e => e.startsWith("Authorization"))[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
let tokenData = jwt_decode(token);
console.log(tokenData);


const store = new Store(initialState);
const popup = new PopupDialog(store);
store.subscribe("isModalUp", [updateModal]);
store.subscribe("loading", [toggleLoader]);
const baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);

const modal = new Modal(store);
const appContainer = document.querySelector("#appContainer");
const deployDialog = document.querySelector("#deployDialog");
const deployDialogForm = document.querySelector("#deployDialog form") as HTMLFormElement;
const deployDialogConfirm = document.querySelector("#btnModalConfirm") as HTMLButtonElement;
deployDialogConfirm.addEventListener("click", e => doForm(e));
const deployDialogCancel = document.querySelector("#btnModalCancel") as HTMLButtonElement;
const searchInp = document.querySelector("#searchInp") as HTMLInputElement;
searchInp.addEventListener("keydown", e => {
	if (e.key == "Backspace" && searchInp.value.length == 1) {
		updateApps();
	}
	updateApps(searchInp.value);
});
const searchBtn = document.querySelector("#searchBtn");
searchBtn.addEventListener("click", () => {
	updateApps(searchInp.value);
});
const deployBtn = document.querySelector("#deployBtn");
deployBtn.addEventListener("click", () => {
});
$("#deployDialog")
	.on("shown.bs.modal", () => store.setState("isModalUp", true))
	.on("hidden.bs.modal", () => store.setState("isModalUp", false));

function init() {
	updateApps();
}

document.addEventListener("keydown", e => {
	switch (e.key) {
		case "Enter":
			if (store.getState("isPopUp")) {
				e.preventDefault();
				popup.confirm();
			} else if (store.getState("isModalUp")) {
				e.preventDefault();
				deployDialogConfirm.click();
			} else if (searchInp == document.activeElement) {
				updateApps(searchInp.value);
			}
			break;
		case "Escape":
			if (store.getState("isPopUp")) {
				popup.cancel();
			} else if (store.getState("isModalUp")) {
				if (modal.up) {
					modal.close();
				} else {
					deployDialogCancel.click();
				}
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
		// let sep = window.location.hostname.split(".");
		// sep.shift();
		// url = window.location.protocol + "//" + sep.join(".") + ":" + port;
		url += window.location.hostname + ":" + port;
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
		case "settings":
			color = "secondary";
			icon = "cogs";
			text = "Settings";
			break;
	}

	return `<button class="btn btn-${color}" data-action=\"${action}\" data-id="${id}" onclick=\"doAction(event)\"><i class="fas fa-${icon} fa-2x"></i><br>${text}</button>`;

}

function dateTemplate(dateString: string): string {
	return new Date(dateString).toLocaleString();
}

function getRunnerIcon(runner: string): string {
	let r = "";
	switch (runner) {
		case "node":
			r = "node";
			break;
		case "web":
			r = "html5";
			break;
		case "npm":
			r = "npm";
			break;
		case "python":
			r = "python";
			break;
		case "flask":
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
                            <span>Hostname:</span><span>${app.hostname == "" ? window.location.hostname + ":" + app.port : app.hostname}</span>
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
                            <span>Runner:</span><span>${getRunnerIcon(app.runner)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Pid:</span><span>${app.pid == -1 ? "offline" : app.pid}</span>
                        </li>
                    </ul>
                </div>
                <div class="card-footer text-right d-flex justify-content-around">
                	${running ? getOpenExternalButton(app.hostname, app.port) : getButton(app.id, "settings")}
                	${running ? getButton(app.id, "kill") : getButton(app.id, "run")}
    				${getButton(app.id, "update")}
    				${getButton(app.id, "remove")}
				</div>
            </div>
        </div>`;
}

function settingsTemplate(app: App): string {
	return `<div class="row"><div class="col"></div><form class="col-md-8">
			<div class="input-group input-group-sm mb-3">
				<div class="input-group-prepend">
					<span class="input-group-text">ID</span>
				</div>
				<input readonly value="${app.id}" type="text" class="form-control" aria-label="Small">
			</div>
			<div class="input-group input-group-sm mb-3">
				<div class="input-group-prepend">
					<span class="input-group-text">Name</span>
				</div>
				<input readonly value="${app.name}" type="text" class="form-control" aria-label="Small">
			</div>
			<div class="input-group input-group-sm mb-3">
				<div class="input-group-prepend">
					<span class="input-group-text">Repository</span>
				</div>
				<input readonly value="${app.repo}" type="text" class="form-control" aria-label="Small">
			</div>
			<div class="input-group mb-3">
				<div class="input-group-prepend">
					<label class="input-group-text" for="runner">Runner</label>
				</div>
				<select class="custom-select" name="runner" id="runnerSettings">
					<option ${app.runner == "node" ? "selected" : ""} value="node">Node</option>
					<option ${app.runner == "npm" ? "selected" : ""} value="npm">Npm</option>
					<option ${app.runner == "web" ? "selected" : ""} value="web">Web</option>
					<option ${app.runner == "python" ? "selected" : ""} value="python">Python</option>
					<option ${app.runner == "flask" ? "selected" : ""} value="flask">Flask</option>
				</select>
			</div>
			<div class="input-group input-group-sm mb-3">
				<div class="input-group-prepend">
					<span class="input-group-text">Hostname</span>
				</div>
				<input value="${app.hostname}" type="text" name="hostname" id="hostnameSettings" class="form-control" aria-label="Small">
			</div>
			<div class="input-group input-group-sm mb-3">
				<div class="input-group-prepend">
					<span class="input-group-text">Port</span>
				</div>
				<input value="${app.port}" type="text" name="port" id="portSettings" class="form-control" aria-label="Small">
			</div>
		</form><div class="col"></div></div>
		<button class="btn btn-success" type="button" data-id="${app.id}" data-action=\"settings\" onclick=\"doModalForm(event)\"">Update</button>`;
}

function updateApps(query: string = "") {
	const url = baseUrl;
	url.pathname = "/api/find";
	url.search = "?app=" + query;
	fetch(url.href).then(j => {
		j.json().then(res => {
			const appsD: App[] = res.deployed != null ? res.deployed : [];
			const apps: App[] = res.running != null ? res.running : [];
			if (query == "") {
				store.setState("deployedApps", appsD);
				store.setState("runningApps", apps);
				console.log(store);
			}
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
	const url = baseUrl;
	const btn = event.target as HTMLButtonElement;
	const action = btn.attributes.getNamedItem("data-action").value;
	const id = btn.attributes.getNamedItem("data-id").value;
	const data = {app: id};
	if (action != "settings") {
		popup.open("Warning", "Are you sure?", () => {
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
				.then(async res => {
					if (res.status == 200) {
						updateApps();
					} else {
						const response = await res.json();
						setTimeout(() => {
							popup.open("Error", response.message.toString().toLocaleUpperCase());
						}, 200);
					}
					store.setState("loading", false);
				})
				.catch(err => {
					store.setState("loading", false);
					console.log(err);
				});
		});
	} else {
		const apps: App[] = store.getState("deployedApps");
		const app = apps.find(a => a.id == id);
		store.setState("currentApp", app);
		modal.open(app.name, settingsTemplate(app));
	}
}


async function doModalForm(event: Event) {
	popup.open("Warning", "Are you sure?", () => {
		event.preventDefault();
		const btn = event.target as HTMLButtonElement;
		const action = btn.attributes.getNamedItem("data-action").value;
		const id = btn.attributes.getNamedItem("data-id").value;
		const app = store.getState("currentApp");
		const url = baseUrl;
		const edits: Edits = {
			"runner": (document.querySelector("#runnerSettings") as HTMLInputElement).value,
			"hostname": (document.querySelector("#hostnameSettings") as HTMLInputElement).value,
			"port": (document.querySelector("#portSettings") as HTMLInputElement).value,
		};
		let settings: Edits = {};
		Object.keys(edits).forEach(key => {
			if (edits[key] != String(app[key])) {
				settings[key] = edits[key];
			}
		});
		const data = {id: id, settings: settings};
		url.pathname = "/api/" + action;
		console.log(action);
		store.setState("loading", true);
		modal.close();
		fetch(url.href, {
			method: "POST", // *GET, POST, PUT, DELETE, etc.
			mode: "cors", // no-cors, cors, *same-origin
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data), // body data type must match "Content-Type" header
		})
			.then(async res => {
				console.log(res);
				if (res.status == 200) {
					updateApps();
				} else {
					const response = await res.json();
					setTimeout(() => {
						popup.open("Error", response.message.toString().toLocaleUpperCase());
					}, 200);
				}
				store.setState("loading", false);
			})
			.catch(async err => {
				console.log(err);
				store.setState("loading", false);
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
		$("#deployDialog").modal("hide");
		fetch(url.href, {
			method: "POST", // *GET, POST, PUT, DELETE, etc.
			mode: "cors", // no-cors, cors, *same-origin
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify(data), // body data type must match "Content-Type" header
		})
			.then(async res => {
				if (res.status == 200) {
					// @ts-ignore
					$("#deployDialog").modal("hide");
					updateApps();
				} else {
					const response = await res.json();
					setTimeout(() => {
						popup.open("Error", response.message.toString().toLocaleUpperCase());
					}, 200);
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
