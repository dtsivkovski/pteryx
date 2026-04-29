/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const (
	cyan = "\033[36m"
	red   = "\033[31m"
	reset = "\033[0m"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pteryx [file]",
	Short: "A file checker that verifies correct file signatures",
	Long: `
           ▓▓▓▓▓                                            
         ▓▒▒▒▒░▒▒█             ▓█▓▓▓▓                       
       ▓▓▓▒▒███▒▓▓           ▓████▓▓███▓███▓▓               
     ▓▓    █▓▓▓█▓▓▓        ▓▓██████████▓██████▓▓▓            
         ▓█   ▓▓▒█     ▓▓▓███████████████▓▓█▓███▓█▓▒        
       ▓▓      ▓▓█   ▓▓██████████████▓▓██     █▓▓██▓█▓      
           ▓▓▓▓▓█▓▓▓▒▓█████████████                  ██▓    
      ▓███▓▓▓███▓██▓████████████                        ▓▓  
     ▓█████▓▓████▓▓█████████▓█                           ██ 
    ▓████████████████████▓██                               █
  ▓█████████▓      ▓████▓█          ==============          
  ▓██████           ▓███▓ █         ==  PTERYX  ==          
 ▓████▓              █▓█ █          ==============          
 ███▓▓                 ██ █                                 
▓███                     ▓                                  
███                       █                                 
██                                                          
 █             

Pteryx is a file checker that verifies correct file signatures.
It helps ensure that files are not intentionally or accidentally misnamed,
made for use in digital forensics and incident response and malware analysis.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}

		// run check on file
		return runFileCheck(args[0])
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pteryx.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
