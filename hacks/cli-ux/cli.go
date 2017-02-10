package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/mgutz/ansi"
	"github.com/spf13/cobra"
)

var (
	gray   = ansi.ColorFunc("black+h")
	cyan   = ansi.ColorFunc("cyan")
	red    = ansi.ColorFunc("red")
	yellow = ansi.ColorFunc("yellow")
	reset  = ansi.ColorFunc("reset")
	bright = ansi.ColorFunc("white+h")
)

func Save(path string, object interface{}) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := json.NewEncoder(file)
		encoder.Encode(object)
	}
	file.Close()
	return err
}

func Load(path string, object interface{}) error {
	file, err := os.Open(path)
	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(object)
	}
	file.Close()
	return err
}

var items map[string]map[string]interface{}

func initKind(kind string) {
	if _, exists := items[kind]; !exists {
		items[kind] = make(map[string]interface{})
	}
}

func AddItem(kind, name string, item interface{}) {
	initKind(kind)
	items[kind][name] = item
}

func RemoveItem(kind, name string) {
	initKind(kind)
	delete(items[kind], name)
}

func ListItems(kind string) map[string]interface{} {
	initKind(kind)
	return items[kind]
}

func PrintItem(kind, name string) {
	initKind(kind)
	PrintFields(items[kind][name].(map[string]interface{}), true)
}

func PrintFields(fields map[string]interface{}, colon bool) {
	longest := 0
	for k, _ := range fields {
		if len(k) > longest {
			longest = len(k) + 2
		}
	}
	for k, v := range fields {
		if colon {
			k = k + ":"
		}
		fmt.Printf(fmt.Sprintf("%%-%ds  %%s\n", longest), k, v)
	}
}

func NewCmd(use string, reqArgs int, action func([]string)) *cobra.Command {
	return &cobra.Command{
		Use: use,
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) < reqArgs {
				cmd.Usage()
				return
			}
			action(args)
		},
	}
}

func NewTable() *tabwriter.Writer {
	return tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
}

func DoneStatus(message string) {
	fmt.Printf("%s... ", message)
	time.Sleep(1 * time.Second)
	fmt.Printf("done\n")
}

func Header(text string) {
	fmt.Println(gray("==="), bright(text))
}

func main() {
	if err := Load("/tmp/cli-ux-state.json", &items); err != nil {
		items = make(map[string]map[string]interface{})
	}

	var rootCmd = &cobra.Command{
		Use: "ssh-cmd [command]",
	}
	rootCmd.AddCommand(
		NewCmd(":create <name>", 1, func(args []string) {
			AddItem("command", args[0], map[string]string{
				"Name":        args[0],
				"Owner":       "progrium",
				"Description": "(none)",
				"Created":     "Mon Jan 1 15:04:05 MST 2016",
				"Updated":     "Mon Jan 2 12:01:00 MST 2016",
				"Size":        "201 MB",
				"Run URL":     "https://cmd.io/run/progrium/" + args[0],
			})
			DoneStatus(fmt.Sprintf("Creating %s", bright(args[0])))
		}))

	rootCmd.AddCommand(
		NewCmd(":destroy <name>", 1, func(args []string) {
			RemoveItem("command", args[0])
			DoneStatus(fmt.Sprintf("Destroying %s", bright(args[0])))
		}))

	rootCmd.AddCommand(
		NewCmd(":ls [flags] [owner]", 0, func(args []string) {
			if len(args) == 1 {
				fmt.Println("gl-deploy")
				return
			}
			Header("Your Commands")
			for name, _ := range ListItems("command") {
				fmt.Println(name)
			}
			fmt.Println()
			Header("Group Commands")
			fmt.Println("gliderlabs/gl-deploy")
			fmt.Println()
			Header("Shared Commands")
			fmt.Println("jpf/rickroll")
			fmt.Println("mattaichison/dune")
		}))

	rootCmd.AddCommand(
		NewCmd(":info <name>", 1, func(args []string) {
			PrintItem("command", args[0])
		}))

	rootCmd.AddCommand(
		NewCmd(":access <name>", 1, func(args []string) {
			if len(args) < 2 {
				Header("Users")
				PrintFields(map[string]interface{}{
					"progrium":      gray("(owner)"),
					"jpf":           gray("(admin)"),
					"mattaitchison": "",
				}, false)
				fmt.Println()
				Header("Groups")
				fmt.Println("gliderlabs")
				fmt.Println()
				Header("Tokens")
				fmt.Println("a17671fb-b2f1–4286–89de-a17671fb")
				fmt.Println("26f471fb-b223–5467–29df-a1767765")
				return
			}
			name := args[0]
			var subCmds = &cobra.Command{Use: "ssh-cmd :access <name>"}
			subCmds.AddCommand(
				NewCmd("grant <subject>", 1, func(args []string) {
					DoneStatus(fmt.Sprintf("Granting %s access to %s", bright(args[0]), bright(name)))
				}),
				NewCmd("revoke <subject>", 1, func(args []string) {
					DoneStatus(fmt.Sprintf("Revoking %s access to %s", bright(args[0]), bright(name)))
				}),
			)
			subCmds.SetArgs(args[1:])
			subCmds.Execute()
		}))

	rootCmd.AddCommand(
		NewCmd(":admins <name>", 1, func(args []string) {
			if len(args) < 2 {
				PrintFields(map[string]interface{}{
					"progrium": gray("(owner)"),
					"jpf":      "",
				}, false)
				return
			}
			name := args[0]
			var subCmds = &cobra.Command{Use: "ssh-cmd :admins <name>"}
			subCmds.AddCommand(
				NewCmd("grant <user>", 1, func(args []string) {
					DoneStatus(fmt.Sprintf("Granting %s admin to %s", bright(args[0]), bright(name)))
				}),
				NewCmd("revoke <user>", 1, func(args []string) {
					DoneStatus(fmt.Sprintf("Revoking %s admin to %s", bright(args[0]), bright(name)))
				}),
			)
			subCmds.SetArgs(args[1:])
			subCmds.Execute()
		}))

	rootCmd.AddCommand(
		NewCmd(":tokens", 0, func(args []string) {
			if len(args) < 1 {
				t := NewTable()
				fmt.Fprintln(t, "TOKEN\tLAST USED\tLAST IP")
				fmt.Fprintln(t, "a17671fb-b2f1–4286–89de-a17671fb\tMon Jan 1 15:04:05 MST 2016\t192.168.0.200")
				fmt.Fprintln(t, "98f471fb-b223–5467–29df-a1767721\t\t")
				fmt.Fprintln(t, "26f471fb-b223–5467–29df-a1767765\tMon Feb 2 15:04:05 MST 2016\t10.0.0.1")
				t.Flush()
				//
				// PrintFields(map[string]interface{}{
				// 	"a17671fb-b2f1–4286–89de-a17671fb": gray(""),
				// 	"26f471fb-b223–5467–29df-a1767765": gray("Mon Jan 1 15:04:05 MST 2016"),
				// }, false)
				return
			}
			var subCmds = &cobra.Command{Use: "ssh-cmd :tokens"}
			subCmds.AddCommand(
				NewCmd("new", 0, func(args []string) {
					fmt.Println("26f471fb-b223–5467–29df-a1767765")
				}),
				NewCmd("rm <token>", 1, func(args []string) {
					DoneStatus(fmt.Sprintf("Removing %s and revoking from all commands", bright(args[0])))
				}),
			)
			subCmds.SetArgs(args[0:])
			subCmds.Execute()
		}))

	rootCmd.AddCommand(
		NewCmd(":env <name>", 1, func(args []string) {
			if len(args) < 2 {
				PrintFields(map[string]interface{}{
					bright("AUTH0_API_TOKEN"):        "eyJhbGciOiJIUzI1NiIsInR5cCI6pXVCJ9.eyJhdWQzBOa3VkZFc5eThiMkhzWGJlT050QXA2TCIsInNjb3BlcyI6eyJ1c2Vyc19hcHBfbWV0YWRhdGEiOnsiYWN0aW9ucyI6WyJyZWFkIiwidXBkYXRlIl19LCJ1c2VycyI6eyJhY3Rpb25zIjpbInJlYWQiLCJ1cGRhdGUiXX19LCJpYXQiOjE0Njg3MjY4MDIsImp0aSI6IjM5MDFhOWFjMGUwZmZjNzczODgzODQ5YjUyODVjMDgwIn0.Ic-TsV0Dm4I2ITOVOPyRCYiy1IiYU",
					bright("GITHUB_ACCESS_TOKEN_CI"): "3b5ace883725b528e4ab88764d75872fb4ae0bff",
					bright("AWS_SECRET_ACCESS_KEY"):  "j0bJFYTLBfTFEZ9McjGJc0COxSJaCkpG71IylxfI",
				}, true)
				return
			}
			name := args[0]
			var subCmds = &cobra.Command{Use: "ssh-cmd :admins <name>"}
			subCmds.AddCommand(
				NewCmd("set <key=value>...", 1, func(args []string) {
					var vars []string
					for _, v := range args {
						parts := strings.SplitN(v, "=", 2)
						vars = append(vars, bright(parts[0]))
					}
					DoneStatus(fmt.Sprintf("Setting %s on %s", strings.Join(vars, ", "), bright(name)))
				}),
				NewCmd("unset <key>...", 1, func(args []string) {
					var vars []string
					for _, v := range args {
						vars = append(vars, bright(v))
					}
					DoneStatus(fmt.Sprintf("Unsetting %s on %s", strings.Join(vars, ", "), bright(name)))
				}),
			)
			subCmds.SetArgs(args[1:])
			subCmds.Execute()
		}))

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

	if err := Save("/tmp/cli-ux-state.json", items); err != nil {
		log.Fatal(err)
	}
}

// ssh cmd.io cmd-create <name> -
// ssh cmd.io cmd-rm <name>
// ssh cmd.io cmd-ls
// ssh cmd.io cmd-info <name>
// ssh cmd.io cmd-access <name> add|rm <user>
// ssh cmd.io cmd-admins <name> add|rm <user>

// ssh cmd.io cmd-tokens new|ls|rm
// ssh cmd.io cmd-env <name> unset|set <key>[=<value>]

// ssh cmd.io cmd-group create|destroy|users <group>

// ssh cmd.io cmd-system version|changes
// ssh cmd.io cmd-account upgrade|cancel|...
// ssh cmd.io cmd-keys add|rm
//
//
// ssh cmd.io :env welcome set FOO=bar
// ssh cmd.io :create foobar
// ssh cmd.io :access
