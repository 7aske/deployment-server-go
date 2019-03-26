import { FindResponse } from "./Responses";
import { App } from "./App";

export function appTemplate(app: App): string {
	return `<div class="card">
                <div class="card-header" id="heading${app.id}">
                    <h2 class="mb-0">
                        <button class="btn btn-link collapsed w-100" type="button" data-toggle="collapse" data-target="#collapse${app.id}" aria-expanded="false" aria-controls="collapse${app.id}">
                            ${app.name}
                        </button>
                    </h2>
                </div>
                <div id="collapse${app.id}" class="collapse" aria-labelledby="heading${app.id}" data-parent="#appContainer">
                    <div class="card-body">
                       ${app.repo}
                    </div>
                </div>
            </div>`;
};

(async function () {
	const json = await fetch("http://localhost:8080/api/find");
	const res: FindResponse = await json.json();
	res.deployed.forEach(a => {
		document.querySelector("#appContainer").innerHTML += appTemplate(a);
	});
})();

