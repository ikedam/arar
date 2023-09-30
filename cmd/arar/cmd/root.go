package cmd

import (
	"context"
	"log"
	"os"
	"os/exec"

	"github.com/ikedam/arar/internal/aws"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "arar",
	Short: "arar performs assume role and run command with that credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		arar := &aws.AssumeRoleAndRun{
			Region:          viper.GetString("region"),
			RoleARN:         viper.GetString("role-name"),
			SessionName:     viper.GetString("role-session-name"),
			UsernameSession: viper.GetBool("username-session"),
		}
		ctx := context.Background()
		session, err := arar.AssumeRole(ctx)
		if err != nil {
			return err
		}
		if len(args) <= 0 {
			return nil
		}
		c := exec.Command(args[0], args[1:]...)
		c.Env = append(
			os.Environ(),
			"AWS_ACCESS_KEY_ID="+session.AccessKeyID,
			"AWS_SECRET_ACCESS_KEY="+session.SecretAccessKey,
			"AWS_SESSION_TOKEN="+session.SessionToken,
		)
		if arar.Region != "" {
			c.Env = append(
				c.Env,
				"AWS_DEFAULT_REGION="+arar.Region,
				"AWS_REGION="+arar.Region,
			)
		}
		return c.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalf("%+v", err)
	}
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	rootCmd.Flags().String("region", "", "AWS region")
	viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))
	rootCmd.Flags().String("role-arn", "", "IAM Role ARN")
	viper.BindPFlag("role-arn", rootCmd.Flags().Lookup("role"))
	rootCmd.Flags().String("role-session-name", "", "session identifier")
	viper.BindPFlag("role-session-name", rootCmd.Flags().Lookup("role-session-name"))
	rootCmd.Flags().BoolP("username-session", "u", false, "use IAM user name as session identifier")
	viper.BindPFlag("username-session", rootCmd.Flags().Lookup("username-session"))

	viper.SetEnvPrefix("AWS")
	viper.BindEnv("role-arn", "role-session-name")
}

// SetVersion sets version of the command
func SetVersion(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
