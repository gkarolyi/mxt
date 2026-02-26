package main

import (
	"errors"
	"os"

	"github.com/gkarolyi/mxt/internal/commands"
	mxtErrors "github.com/gkarolyi/mxt/internal/errors"
	"github.com/gkarolyi/mxt/internal/ui"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

const version = "1.0.0"

var rootCmd = &cobra.Command{
	Use:   "mxt",
	Short: "Tmux Worktree Session Manager",
	Long: `                       _
  _ __ ___  _   ___  _| |_ _ __ ___  ___
 | '_ ` + "`" + ` _ \| | | \ \/ / __| '__/ _ \/ _ \
 | | | | | | |_| |>  <| |_| | |  __/  __/
 |_| |_| |_|\__,_/_/\_\__|_|  \___|\___|

  Tmux Worktree Session Manager v` + version + `

A tool for managing git worktrees paired with tmux sessions.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion, _ := cmd.Flags().GetBool("version")
		if showVersion {
			commands.VersionCommand(version)
			return
		}
		commands.HelpCommand(version)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version number",
	Run: func(cmd *cobra.Command, args []string) {
		commands.VersionCommand(version)
	},
}

var initCmd = &cobra.Command{
	Use:   "init [--local] [--reinit]",
	Short: "Set up configuration",
	Long:  "Create global config (~/.mxt/config) or project config (.mxt in repo root)",
	Run: func(cmd *cobra.Command, args []string) {
		local, _ := cmd.Flags().GetBool("local")
		reinit, _ := cmd.Flags().GetBool("reinit")
		if err := commands.InitCommand(local, reinit); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := commands.ConfigCommand(); err != nil {
			var notFound mxtErrors.ErrConfigNotFound
			if errors.As(err, &notFound) {
				os.Exit(1)
			}
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var newCmd = &cobra.Command{
	Use:   "new [branch-name]",
	Short: "Create worktree + tmux session",
	Run: func(cmd *cobra.Command, args []string) {
		branchName := ""
		if len(args) == 0 {
			if term.IsTerminal(int(os.Stdin.Fd())) {
				prompted, err := commands.PromptBranchName(os.Stdin, os.Stdout)
				if err != nil {
					ui.Error(err.Error())
					os.Exit(1)
				}
				branchName = prompted
			} else {
				ui.Error("Usage: mxt new [branch-name] [--from <base-branch>] [--run claude|codex] [--bg]")
				os.Exit(1)
			}
		} else {
			branchName = args[0]
		}
		fromBranch, _ := cmd.Flags().GetString("from")
		runCmd, _ := cmd.Flags().GetString("run")
		bg, _ := cmd.Flags().GetBool("bg")

		if err := commands.NewCommand(branchName, fromBranch, runCmd, bg); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List worktrees, diff stats, session status",
	Run: func(cmd *cobra.Command, args []string) {
		if err := commands.ListCommand(); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var deleteCmd = &cobra.Command{
	Use:     "delete <branch-name>",
	Aliases: []string{"rm"},
	Short:   "Delete worktree and branch",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			ui.Error("Usage: mxt delete <branch-name> [--force|-f]")
			os.Exit(1)
		}
		force, _ := cmd.Flags().GetBool("force")
		if err := commands.DeleteCommand(args[0], force); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var sessionsCmd = &cobra.Command{
	Use:   "sessions <action> <branch-name>",
	Short: "Manage tmux session for a worktree",
	Long: `Actions:
  open   <branch> [--run cmd]   Create session & open terminal
  close  <branch>               Kill tmux session
  relaunch <branch> [--run cmd] Close + reopen session
  attach <branch> [dev|agent]   Attach to session (optionally select window)`,
	Args: cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			ui.Error("Usage: mxt sessions <open|close|relaunch|attach> <branch> [--run claude|codex] [--bg]")
			os.Exit(1)
		}
		action := args[0]
		branchName := ""
		if len(args) > 1 {
			branchName = args[1]
		}
		if branchName == "" {
			switch action {
			case "open", "launch", "start":
				ui.Error("Usage: mxt sessions open <branch> [--run claude|codex] [--bg]")
			case "close", "kill", "stop":
				ui.Error("Usage: mxt sessions close <branch>")
			case "relaunch", "restart":
				ui.Error("Usage: mxt sessions relaunch <branch> [--run claude|codex] [--bg]")
			case "attach":
				ui.Error("Usage: mxt sessions attach <branch> [dev|agent]")
			default:
				ui.Error("Usage: mxt sessions <open|close|relaunch|attach> <branch> [--run claude|codex] [--bg]")
			}
			os.Exit(1)
		}

		// For attach command, third argument can be window name
		windowName := ""
		if action == "attach" && len(args) > 2 {
			windowName = args[2]
		}

		runCmd, _ := cmd.Flags().GetString("run")
		bg, _ := cmd.Flags().GetBool("bg")

		// Pass windowName as runCmd for attach action (reusing parameter)
		if action == "attach" {
			runCmd = windowName
		}

		if err := commands.SessionsCommand(action, branchName, runCmd, bg); err != nil {
			ui.Error(err.Error())
			os.Exit(1)
		}
	},
}

var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "Display help information",
	Run: func(cmd *cobra.Command, args []string) {
		commands.HelpCommand(version)
	},
}

func init() {
	// Add flags for init command
	initCmd.Flags().BoolP("local", "l", false, "Create project config (.mxt in repo root)")
	initCmd.Flags().Bool("reinit", false, "Overwrite existing config without prompting")

	// Add flags for new command
	newCmd.Flags().String("from", "", "Base branch (default: main/master)")
	newCmd.Flags().String("run", "", "Auto-run command in agent window (claude|codex)")
	newCmd.Flags().Bool("bg", false, "Create session without opening terminal")

	// Add flags for delete command
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")

	// Add flags for sessions command
	sessionsCmd.Flags().String("run", "", "Auto-run command in agent window (claude|codex)")
	sessionsCmd.Flags().Bool("bg", false, "Create session without opening terminal")

	// Add subcommands to root
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(sessionsCmd)
	rootCmd.AddCommand(helpCmd)
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		commands.HelpCommand(version)
	})

	// Customize version flag
	rootCmd.Flags().BoolP("version", "v", false, "Display version number")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}
