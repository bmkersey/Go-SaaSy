/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Go-SaaSy",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
	rootCmd.AddCommand(createUserCmd, createOrgCmd, resetPasswordCmd)

	createUserCmd.Flags().StringVar(&userEmail, "email", "", "User's email")
	createUserCmd.Flags().StringVar(&userPassword, "password", "", "User's password")

	createOrgCmd.Flags().StringVar(&orgName, "name", "", "Organization name")
	createOrgCmd.Flags().StringVar(&ownerID, "owner-id", "", "UUID of the user who will own this org")

	resetPasswordCmd.Flags().StringVar(&resetEmail, "email", "", "User email to reset password for")
	resetPasswordCmd.Flags().BoolVar(&sendEmail, "send", false, "Send the password reset email")
}
