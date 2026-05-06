/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	magenta = "\033[35m"
	cyan    = "\033[36m"
	red     = "\033[31m"
	reset   = "\033[0m"
)

var version = "0.0.3"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pteryx",
	Short:   "A forensic file inspection tool",
	Version: version,
	Long: `

             █▓▓        
             ▒▒▒▒▒▒     
            ▒▓     ▓▓▓  
           ▒▒           
          ▒▒      ██████╗ ████████╗███████╗██████╗ ██╗   ██╗██╗  ██╗
       ▓▓▒▒▓      ██╔══██╗╚══██╔══╝██╔════╝██╔══██╗╚██╗ ██╔╝╚██╗██╔╝
    ▓ ▓▒▓▒▓▓      ██████╔╝   ██║   █████╗  ██████╔╝ ╚████╔╝  ╚███╔╝ 
    ▓ ▓▒ ▓▓▒▒     ██╔═══╝    ██║   ██╔══╝  ██╔══██╗  ╚██╔╝   ██╔██╗ 
    ▓▓▓    ▓▓     ██║        ██║   ███████╗██║  ██║   ██║   ██╔╝ ██╗
    ▓▓▓     ██    ╚═╝        ╚═╝   ╚══════╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝      
      ▒█                           

Pteryx is a forensic file inspection tool.
It verifies file signatures and can create hash baselines for later comparison,
made for use in digital forensics, incident response, and malware analysis.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
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
	rootCmd.AddCommand(sigCmd)
	rootCmd.AddCommand(hashCmd)
}
