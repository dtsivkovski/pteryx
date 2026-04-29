/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)



// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pteryx",
	Short: "A file checker that verifies correct file signatures",
	Long: `
           ▓▓▓▓▓                                            
         ▓▒▒▒▒░▒▒█             ▓█▓▓▓▓                       
       ▓▓▓▒▒███▒▓▓           ▓████▓▓███▓███▓▓               
     ▓    █▓▓▓█▓▓▓        ▓▓██████████▓██████▓▓▓            
         ▓█   ▓▓▒█     ▓▓▓███████████████▓▓█▓███▓█▓▒        
       ▓▓      ▓▓█   ▓▓██████████████▓▓██     █▓▓██▓█▓      
           ▓▓▓▓▓█▓▓▓▒▓█████████████                  ██▓    
      ▓███▓▓▓███▓██▓████████████                        ▓▓  
     ▓█████▓▓████▓▓█████████▓█                           ██ 
    ▓████████████████████▓██                               █
  ▓█████████▓      ▓████▓█                                  
  ▓██████           ▓███▓ █                                 
 ▓████▓              █▓█ █                                  
 ███▓▓                 ██ █                                 
▓███                     ▓                                  
███                       █                                 
██                                                          
 █                                                          
Pteryx is a file checker that verifies correct file signatures.
It helps ensure that files are not intentionally or accidentally misnamed,
made for use in digital forensics and incident response and malware analysis.
`,
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pteryx.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


