package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/sarcasticadmin/sshcb/builder"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sshcb",
	Short: "sshcb can easy way to build ssh_configs",
	Long: `Connect to your environment quickly and easily
	       by querying a cloud api and building an ssh_config`,
	Run: func(cmd *cobra.Command, args []string) {
		rawtags, _ := cmd.Flags().GetStringSlice("tags")
		tags := make(map[string]string)
		for _, value := range rawtags {
			tmp := strings.Split(value, ":")
			tags[tmp[0]] = tmp[1]
		}

		region, _ := cmd.Flags().GetString("region")
		profile, _ := cmd.Flags().GetString("profile")
		username, _ := cmd.Flags().GetString("username")
		bastionhost, _ := cmd.Flags().GetString("bastionhost")
		outputfile, _ := cmd.Flags().GetString("output")
		session := builder.GetSession(profile, region)
		resp := builder.GetReservs(tags, session)
		myConfig := builder.SSHConfigOptions{
			Username:    username,
			Filepath:    outputfile,
			BastionHost: bastionhost}
		instances := builder.BuildInstanceList(resp.Reservations)
		builder.WriteSSHConfig(instances, myConfig)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SSHcb",
	Long:  `All software has versions. This is SSHcb`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v0.0")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize()
	rootCmd.PersistentFlags().StringP("region", "r", "us-east-1", "Datacenter region")
	rootCmd.PersistentFlags().StringP("username", "u", "ec2-user", "SSH Username")
	rootCmd.PersistentFlags().StringP("output", "o", "./config", "Output Location of SSH Config")
	rootCmd.PersistentFlags().StringP("profile", "p", "default", "AWS profile to use from ~/.aws/credentials")
	rootCmd.PersistentFlags().StringP("bastionhost", "b", "", "bastion IP or hostname into AWS")
	rootCmd.PersistentFlags().StringSlice("tags", []string{}, "instance tags AWS in the form of key:value")
	rootCmd.AddCommand(versionCmd)
}
