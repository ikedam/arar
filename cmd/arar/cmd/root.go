package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/ikedam/arar/internal/aws"
	"github.com/ikedam/arar/internal/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:           "arar",
	Short:         "arar performs assume role and run command with that credentials",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if viper.GetBool("verbose") {
			log.SetLevelByName("Debug")
		}
		arar := &aws.AssumeRoleAndRun{
			Region:          viper.GetString("region"),
			RoleARN:         viper.GetString("role_arn"),
			SessionName:     viper.GetString("role_session_name"),
			UsernameSession: viper.GetBool("username_session"),
			SerialNumber:    viper.GetString("serial_number"),
			DurationSeconds: viper.GetInt32("duration_seconds"),
		}
		ctx := log.CtxWithLogger(context.Background())
		session, err := arar.AssumeRole(ctx)
		if err != nil {
			return err
		}
		if len(args) <= 0 {
			return nil
		}
		c := exec.Command(args[0], args[1:]...)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
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
		log.Debug(
			ctx,
			"call command",
			log.WithStrings("command", args),
		)
		return c.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(context.Background(), "command failed", log.WithError(err))
	}
}

func init() {
	rootCmd.Flags().BoolP("verbose", "v", false, "verbose output")
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	rootCmd.Flags().String("region", "", "AWS region")
	viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))
	rootCmd.Flags().String("role-arn", "", "IAM Role ARN")
	viper.BindPFlag("role_arn", rootCmd.Flags().Lookup("role-arn"))
	rootCmd.Flags().String("role-session-name", "", "session identifier")
	viper.BindPFlag("role_session_name", rootCmd.Flags().Lookup("role-session-name"))
	rootCmd.Flags().BoolP("username-session", "u", false, "use IAM user name as session identifier")
	viper.BindPFlag("username_session", rootCmd.Flags().Lookup("username-session"))
	rootCmd.Flags().String("serial-number", "", "MFA device number or virtual MFA device ARN")
	viper.BindPFlag("serial_number", rootCmd.Flags().Lookup("serial-number"))
	rootCmd.Flags().String("duration-seconds", "", "expiration time for role session")
	viper.BindPFlag("duration_seconds", rootCmd.Flags().Lookup("duration-seconds"))

	viper.SetEnvPrefix("aws")
	viper.BindEnv("role_arn")
	viper.BindEnv("role_session_name")
}

// SetVersion sets version of the command
func SetVersion(version string) {
	rootCmd.Version = version
	rootCmd.SetVersionTemplate("{{.Version}}\n")
}
