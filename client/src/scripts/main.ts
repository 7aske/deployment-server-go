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
const baseUrl = new URL(window.location.protocol + "//" + window.location.hostname + ":" + window.location.port);
const token = document.cookie.split("; ").filter(e => e.startsWith("Authorization"))[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
let tokenData = jwt_decode(token);

function init() {
	let url = baseUrl;
	url.pathname = "/api/find";
	fetch(url.href).then(j => {
		j.json().then(res => {
			console.log(res);
			const appsD: App[] = res.deployed != null ? res.deployed : [];
			const apps: App[] = res.running != null ? res.running : [];
			appsD.forEach(a => {
				if (apps.filter(app => a.id == app.id).length == 1) {
					document.querySelector("#appContainer").innerHTML += appTemplate(apps.find(app => a.id == app.id), true);
				} else {
					document.querySelector("#appContainer").innerHTML += appTemplate(a, false);
				}
			});
		}).catch(err => console.log(err));
	}).catch(err => console.log(err));
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
                <h3 class="mb-2" style="cursor: pointer;" data-toggle="collapse" data-target="#collapse${app.id}" aria-expanded="false" aria-controls="collapse${app.id}">
					${app.name}
                </h3>
                <h6 class="mb-0 text-muted">
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
                <div class="card-footer text-right">
                	${running ? getButton(app.id, "kill") : getButton(app.id, "run")}
    				${getButton(app.id, "update")}
    				${getButton(app.id, "remove")}
				</div>
            </div>
        </div>`;
}

function doAction(event: Event) {
	let url = baseUrl;
	const btn = event.target as HTMLButtonElement;
	const action = btn.attributes.getNamedItem("data-action").value;
	const id = btn.attributes.getNamedItem("data-id").value;
	let data = {app: id};
	url.pathname = `/api/${action}`;
	fetch(url.href, {
		method: "POST", // *GET, POST, PUT, DELETE, etc.
		mode: "cors", // no-cors, cors, *same-origin
		headers: {
			"Content-Type": "application/json",
		},
		body: JSON.stringify(data), // body data type must match "Content-Type" header
	})
		.then(res => res.status == 200 ? location.reload() : null)
		.catch(err => console.log(err));
}

init();
