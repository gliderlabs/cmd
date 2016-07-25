package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/progrium/go-shell"
	"github.com/spf13/cobra"
)

var allowedUsers = []string{"progrium", "mdiebolt", "mattaitchison"}

func main() {
	cmd := filepath.Base(os.Args[0])
	if cmd == "auth" {
		os.Args = append([]string{os.Args[0], cmd}, os.Args[1:]...)
		rootCmd.AddCommand(cmdAuth)
	}
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "cmd",
	Short: "cmd powers cmd.io",
	Run: func(cmd *cobra.Command, args []string) {
		if !allowed(os.Getenv("USER")) {
			fmt.Println("Coming soon...")
			return
		}
		if len(args) == 0 {
			fmt.Println("No command specified.")
			os.Exit(1)
		}
		subcmds := strings.Split(args[0], ":")
		if len(subcmds) > 1 {
			switch subcmds[1] {
			case "config":
				configSubcmd(os.Getenv("USER"), subcmds[0], args[1:])
			default:
				fmt.Println("Unknown subcommand.")
				os.Exit(1)
			}
			return
		}
		image := fmt.Sprintf("%s/cmd-%s", os.Getenv("USER"), args[0])
		pull := exec.Command("/usr/bin/docker", "pull", image)
		pull.Stderr = os.Stderr
		if err := pull.Run(); err != nil {
			os.Exit(1)
		}
		envFile := ensureConfig(os.Getenv("USER"), args[0])
		docker := exec.Command("/usr/bin/docker",
			append([]string{"run", "--rm", "--env-file", envFile, image}, args[1:]...)...)
		docker.Stdout = os.Stdout
		docker.Stderr = os.Stderr
		if err := docker.Run(); err != nil {
			os.Exit(1)
		}
	},
}

func configSubcmd(user, cmd string, args []string) {
	basePath := "/config/" + user + "/" + cmd
	if len(args) == 0 {
		f, err := os.Open(basePath + ".env")
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, f)
		return
	}
	switch args[0] {
	case "set":
		for _, kvp := range args[1:] {
			parts := strings.SplitN(kvp, "=", 2)
			if len(parts) < 2 {
				continue
			}
			err := ioutil.WriteFile(basePath+"/"+strings.Trim(parts[0], "/. "), []byte(parts[1]), 0644)
			if err != nil {
				log.Fatal(err)
			}
			renderConfig(user, cmd)
		}
	case "unset":
		for _, key := range args[1:] {
			os.Remove(basePath + "/" + strings.Trim(key, "/. "))
		}
		renderConfig(user, cmd)
	default:
		fmt.Println("Unknown subcommand.")
		os.Exit(1)
	}
	return
}

func renderConfig(user, cmd string) {
	basePath := "/config/" + user + "/" + cmd
	config := make(map[string]string)
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		if !file.IsDir() {
			b, err := ioutil.ReadFile(basePath + "/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}
			config[file.Name()] = string(b)
		}
	}
	output := ""
	for k, v := range config {
		output = output + k + "=" + v + "\n"
	}
	err = ioutil.WriteFile(basePath+".env", []byte(output), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func ensureConfig(user, cmd string) string {
	envFile := "/config/" + user + "/" + cmd + ".env"
	shell.Run("mkdir -p /config/" + user + "/" + cmd)
	shell.Run("touch " + envFile)
	return envFile
}

var cmdAuth = &cobra.Command{
	Use: "auth <user> <key>",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			os.Exit(2)
			return
		}
		if os.Getenv("CMD_NOAUTH") != "" {
			log.Println("Warning: CMD_NOAUTH is set allowing all SSH connections")
			return
		}
		if !githubKeyAuth(args[0], args[1]) {
			log.Println("auth[ssh]: not allowing", args[0])
			os.Exit(1)
		}
		log.Println("auth[ssh]: allowing", args[0])
	},
}

func githubKeyAuth(user, key string) bool {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", fmt.Sprintf("https://github.com/%s.keys", user), nil)
	resp, _ := client.Do(req)
	if resp.StatusCode != 200 {
		return false
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), key)
}

func allowed(user string) bool {
	for _, u := range allowedUsers {
		if u == user {
			return true
		}
	}
	return false
}
