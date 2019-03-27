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

window.addEventListener("keypress", e => {
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
	const link = location.hostname + ":" + port;
	return `<button class="btn btn-secondary" onclick="window.open('${link}', '_blank')"><i class="fas fa-external-link-alt fa-2x"></i><br>Open</button>`;

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
                	${running ? getOpenExternalButton("", app.port) : ""}
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
			console.log(err)
		});
}

function doForm(event: Event) {
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
}

init();
