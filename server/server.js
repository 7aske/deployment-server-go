const http = require("http");
const fs = require("fs");
const path = require("path");
const mimeTypes = {
    ".html": "text/html",
    ".js": "text/javascript",
    ".css": "text/css",
    ".json": "application/json",
    ".png": "image/png",
    ".jpg": "image/jpg",
    ".gif": "image/gif",
    ".wav": "audio/wav",
    ".mp4": "video/mp4",
    ".woff": "application/font-woff",
    ".ttf": "application/font-ttf",
    ".eot": "application/vnd.ms-fontobject",
    ".otf": "application/font-otf",
    ".svg": "application/image/svg+xml"
};
console.log("Starting server on port " + process.env.PORT);
http.createServer((request, response) => {
    let filePath = "." + request.url;
    if (filePath == "./") {
        filePath = "./index.html";
    }
    let extname = String(path.extname(filePath)).toLowerCase();
    let contentType = mimeTypes[extname] || "text/plain";
    fs.readFile(filePath, (error, content) => {
        if (error) {
            if(error.code == "ENOENT") {
                response.writeHead(404, { "Content-Type": "application/json" });
                response.end("( ͡° ʖ̯ ͡°) 404 NOT FOUND");
            }
            else {
                response.writeHead(500, { "Content-Type": "application/json" });
                response.end("( ͠° ͟ʖ ͡°) 500 INTERNAL SERVER ERROR");
            }
        }
        else {
            response.writeHead(200, { "Content-Type": contentType });
            response.end(content);
        }
    });

}).listen(process.env.PORT);