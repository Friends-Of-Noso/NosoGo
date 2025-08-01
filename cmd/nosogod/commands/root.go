package commands

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"

	ver "github.com/Friends-Of-Noso/NosoGo/version"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Version: ver.Version,
	Use:     fmt.Sprintf("%sd", ver.Name),
	Short:   "The node for the NOSO crypto coin",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// log.Debug("root.init")

	// Custom Usage Function
	rootCmd.SetUsageFunc(rootUsageFunc)

	// rootCmd.SetVersionTemplate(fmt.Sprintf("%sd\n", ver.Title))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func rootUsageFunc(cmd *cobra.Command) error {
	fmt.Print("\033[1mUSAGE\033[0m")
	if cmd.Runnable() {
		fmt.Printf("\n  %s", cmd.UseLine())
	}
	if cmd.HasAvailableSubCommands() {
		fmt.Printf("\n  %s [command]", cmd.CommandPath())
	}
	if len(cmd.Aliases) > 0 {
		fmt.Printf("\n\n\033[1mALIASES\033[0m\n")
		fmt.Printf("  %s", cmd.NameAndAliases())
	}
	if cmd.HasExample() {
		fmt.Printf("\n\n\033[1mEXAMPLES\033[0m\n")
		fmt.Printf("%s", cmd.Example)
	}
	if cmd.HasAvailableSubCommands() {
		cmds := cmd.Commands()
		if len(cmd.Groups()) == 0 {
			fmt.Printf("\n\n\033[1mAVAILABLE COMMANDS\033[0m")
			for _, subcmd := range cmds {
				if subcmd.IsAvailableCommand() || subcmd.Name() == "help" {
					fmt.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
				}
			}
		} else {
			for _, group := range cmd.Groups() {
				fmt.Printf("\n\n%s", group.Title)
				for _, subcmd := range cmds {
					if subcmd.GroupID == group.ID && (subcmd.IsAvailableCommand() || subcmd.Name() == "help") {
						fmt.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
			if !cmd.AllChildCommandsHaveGroup() {
				fmt.Printf("\n\n\033[1mADDITIONAL COMMANDS\033[0m")
				for _, subcmd := range cmds {
					if subcmd.GroupID == "" && (subcmd.IsAvailableCommand() || subcmd.Name() == "help") {
						fmt.Printf("\n  %s %s", rpad(subcmd.Name(), subcmd.NamePadding()), subcmd.Short)
					}
				}
			}
		}
	}
	if cmd.HasAvailableLocalFlags() {
		fmt.Printf("\n\n\033[1mFLAGS\033[0m\n")
		fmt.Print(trimRightSpace(cmd.LocalFlags().FlagUsages()))
	}
	if cmd.HasAvailableInheritedFlags() {
		fmt.Printf("\n\n\033[1mGLOBAL FLAGS\033[0m\n")
		fmt.Print(trimRightSpace(cmd.InheritedFlags().FlagUsages()))
	}
	if cmd.HasHelpSubCommands() {
		fmt.Printf("\n\n\033[1mADDITIONAL HELP TOPICS\033[0m")
		for _, subcmd := range cmd.Commands() {
			if subcmd.IsAdditionalHelpTopicCommand() {
				fmt.Printf("\n  %s %s", rpad(subcmd.CommandPath(), subcmd.CommandPathPadding()), subcmd.Short)
			}
		}
	}

	if cmd.HasAvailableSubCommands() {
		fmt.Printf("\n\nUse \"%s [command] --help\" for more information about a command.", cmd.CommandPath())
	}
	fmt.Println()
	return nil
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}

func lpad(s string, padding int) string {
	return strings.Repeat(" ", padding) + s
}
