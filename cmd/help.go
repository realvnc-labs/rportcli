package cmd

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

More help online: https://cli.rport.io`

var environmentVariables = `

Environment Variables:

The reportcli support use of the following environment variables. 

CONFIG_PATH               specify config.json file location
CONN_TIMEOUT_SEC          set connection timeout
SESSION_VALIDITY_SECONDS  initial lifetime of interactive command session

The following environment variables are used during authentication. See Server Authentication below for more info on use.

RPORT_API_URL       the URL of the rportd server 
RPORT_API_USER      the user name for the authentication
RPORT_API_PASSWORD  the password to be used (do not set if using an API token)'
RPORT_API_TOKEN     your API token (will not work when RPORT_API_PASSWORD is set)
`

var serverAuthentication = `

Server Authentication:

The rportcli requires authentication with the rportd server. The CLI supports three authentication methods.

1) The most straightforward method is to use the RPORT_API_USER and RPORT_API_TOKEN with an API token. This will authenticate directly with the server and bypass any 2fa. This is the best option for unattended CLI use.

2) Alternatively, RPORT_API_USER and RPORT_API_PASSWORD can be used (will not work with 2fa unless using the init command)'. If 2fa is enabled then the authentication will fail and either method 1) or 3) must be used.

3) Using the init command to set an authentication token in the config.json file. If 2fa is enabled then the user's authentication token will be saved in a config.json (see readme for more information) file. If 2fa is enabled then the user will be taken through the 2fa flow and the final authentication token will be saved in config.json.

`

var serverAuthenticationRefer = `

Note:
All commands must authenticate with the rportd server. Please run "rportcli --help" for more information.
`
