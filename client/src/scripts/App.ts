export interface App {
	id: string;
	repo: string;
	name: string;
	root: string;
	port: number;
	hostname: string;
	deployed: Date;
	lastUpdated: Date;
	lastRun: Date;
	uptime: Date;
	runner: string;
	pid?: number;
}

// type AppJSON struct {
// 	Id          string    `json:"id"`
// 	Repo        string    `json:"repo"`
// 	Name        string    `json:"name"`
// 	Root        string    `json:"root"`
// 	Port        int       `json:"port"`
// 	Hostname    string    `json:"hostname"`
// 	Deployed    time.Time `json:"deployed"`
// 	LastUpdated time.Time `json:"last_updated"`
// 	LastRun     time.Time `json:"last_run"`
// 	Uptime      time.Time `json:"uptime"`
// 	Runner      string    `json:"runner"`
// }
