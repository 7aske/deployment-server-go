package utils

import (
	"fmt"
	url2 "net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func RenderLoginPage() []byte {
	return []byte(`<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8" />
		<meta
			name="viewport"
			content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0"
		/>
		<meta http-equiv="X-UA-Compatible" content="ie=edge" />
		<title>Deployment Server</title>
		<style>
			* {
				font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif,
					"Apple Color Emoji", "Segoe UI Emoji", "Segoe UI Symbol", "Noto Color Emoji";
			}
			body,
			html {
				width: 100%;
				overflow-x: hidden;
				background-color: #222222;
				color: #ffffff;
			}
			input {
				margin-left: -8px;
				border-radius: 8px;
				padding: 10px;
				border: 2px solid #666666;
				font-size: 24px;
				margin-bottom: 10px;
			}
			button {
				font-size: 24px;
				border: 2px solid #666666;
				border-radius: 8px;
				padding: 10px 25px;
				background-color: gold;
			}
		</style>
	</head>
	<body style="text-align: center;">
		<h1>Admin Login</h1>
		<form method="POST" action="/auth">
			<input type="text" name="username" placeholder="Username" /><br />
			<input type="password" name="password" placeholder="Password" /><br />
			<button type="submit">Login</button>
		</form>
	</body>
</html>`)
}

func byteCountBinary(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(b)/float64(div), "KMGTPE"[exp])
}

func ParseArgs(q string) (string, bool) {
	index := Contains(q, &os.Args)
	if index == -1 {
		return "", false
	}
	if len(os.Args) == index+1 {
		return "", false
	}
	return os.Args[index+1], true
}
func Contains(q string, s *[]string) int {
	for i, str := range *s {
		if str == q {
			return i
		}
	}
	return -1
}

func ContainsFile(q string, dir *[]os.FileInfo) bool {
	for _, file := range *dir {
		if file.Name() == q {
			return true
		}
	}
	return false
}
func PrintHelp() {
	fmt.Println("usage: [...options] [...flags]")
	fmt.Println()
	fmt.Println("-i		interactive shell access")
}

func GetNameFromRepo(repo string) string {
	reg, _ := regexp.Compile("([^/]+$)")
	return string(reg.Find([]byte(repo)))
}

func PathExits(pth string) bool {
	_, err := os.Stat(pth)
	return err == nil
}

func MakeDirIfNotExist(pth string) {
	if !PathExits(pth) {
		_ = os.MkdirAll(pth, 0775)
	}
}

func MakeFileIfNotExist(pth string) {
	if !PathExits(pth) {
		_ = os.MkdirAll(path.Dir(pth), 0775)
		fp, _ := os.Create(pth)
		_ = fp.Close()
		gid := os.Getgid()
		uid := os.Getuid()
		if gid == 0 {
			if id, err := strconv.Atoi(os.Getenv("SUDO_GID")); err == nil {
				gid = id
			}
		}
		if uid == 0 {
			if id, err := strconv.Atoi(os.Getenv("SUDO_UID")); err == nil {
				uid = id
			}
		}
		_ = os.Chown(pth, uid, gid)
	}
}

func GetAbsDir(pth string) string {
	if path.IsAbs(pth) {
		source, err := os.Readlink(pth)
		if err != nil {
			return os.Args[0]
		} else {
			return source
		}
	} else {
		out, err := exec.Command("which", pth).Output()
		if err != nil {
			return ""
		}
		link := strings.TrimRight(string(out), "\n")
		source, err := os.Readlink(link)
		if err != nil {
			return link
		} else {
			return source
		}
	}
}

func VerifyCcont() bool {
	fmt.Println("verifying ccont")
	return exec.Command("ccont", "-l").Run() == nil
}

func AppendProtocol(repo string) (string, error) {
	url, err := url2.Parse(repo)
	if err != nil {
		return "", err
	}
	if url.Scheme == "" {
		url.Scheme = "https"
	}
	return url.String(), nil
}
