export interface App {
	id: string;
	repo: string;
	name: string;
	root: string;
	port: number;
	hostname: string;
	deployed: string;
	last_updated: string;
	last_run: string;
	uptime: number;
	runner: string;
	pid?: number;
}

export interface FindResponse {
	running: App[];
	deployed: App[];
}
function dateTemplate(dateString: string): string {
	console.log(dateString);
	const date = new Date(dateString);
	return date.toDateString();
	// return ``;
}

function appTemplate(app: App, running: boolean): string {
	return ` <div class="card">
            <div class="card-header" id="heading${app.id}">
				<span class="float-right ${running ? "online text-success" : "offline text-danger"}">${running ? "Online &bull;" : "Offline"}</span>
                <h2 class="mb-0" style="cursor: pointer;" data-toggle="collapse" data-target="#collapse${app.id}" aria-expanded="false" aria-controls="collapse${app.id}">
					${app.name}
                </h2>
                <h5 class="mb-0 text-muted">
					${app.repo}
				</h5>
            </div>
            <div id="collapse${app.id}" class="collapse" aria-labelledby="heading${app.id}" data-parent="#appContainer">
                <div class="card-body row">
                    <ul class="list-group list-group-flush col-lg-6 col-md-12">
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
                            <span>Port:</span><span>${app.port}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Hostname:</span><span>${app.hostname}</span>
                        </li>
                    </ul>
                    <ul class="list-group list-group-flush col-lg-6 col-md-12">
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Deployed:</span><span>${app.deployed}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>LastUpdated:</span><span>${dateTemplate(app.last_updated)}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>LastRun:</span><span>${app.last_run}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Uptime:</span><span>${running ? app.uptime : "offline"}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Runner:</span><span>${app.runner}</span>
                        </li>
                        <li class="list-group-item d-flex justify-content-between">
                            <span>Pid:</span><span>${app.pid}</span>
                        </li>
                    </ul>
                </div>
            </div>
        </div>`;
};
const url = new URL(window.location.href);
const token = document.cookie.split("; ").filter(e => e.startsWith("Authorization"))[0].split("Bearer ")[1].replace("\"", " ");
// @ts-ignore
let tokenData = jwt_decode(token);
console.log(tokenData);

fetch(url.href + "/api/find").then(j => {
	j.json().then(res => {
		console.log(res);
		const appsD: App[] = res.deployed;
		const apps: App[] = res.running != null ? res.running : [];
		appsD.forEach(a => {
			const id = a.id;
			if (apps.filter(app => id == app.id).length == 1) {
				document.querySelector("#appContainer").innerHTML += appTemplate(apps.find(app => id == app.id), true);
			} else {
				document.querySelector("#appContainer").innerHTML += appTemplate(a, false);
			}
		});
	}).catch(err => console.log(err));
}).catch(err => console.log(err));

