package cmd

import "github.com/spf13/cobra"

var sigCmd = &cobra.Command{
	Use:   "sig <file-or-directory>",
	Short: "Check file extensions against their file signatures",
	Long: `Check whether files match the extensions they claim to be.

Use -d for a directory and -r with -d (or -rd) to scan recursively.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		allowDirectory, err := cmd.Flags().GetBool("directory")
		if err != nil {
			return err
		}

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		return runPathCheck(args[0], allowDirectory, recursive)
	},
}

func init() {
	sigCmd.Flags().BoolP("directory", "d", false, "check files in a directory")
	sigCmd.Flags().BoolP("recursive", "r", false, "recursively check directories")
}
