package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/mitchellh/colorstring"
)

// CLI is the command line interface object.
type CLI struct {
	inStream             io.Reader
	outStream, errStream io.Writer
}

const (
	// EnvDebug is environmental var to handle debug mode
	EnvDebug = "GO_DEBUG"
)

// Exit codes are values representing an exit code for a error type.
const (
	ExitCodeOK int = 0

	// Errors start at 10
	ExitCodeError = 10 + iota
	ExitCodeParseFlagsError

	ExitCodeGitHubUsersDoNotExistError
	ExitCodeSlackWebhookUrlDoNotExistError
	ExitCodeSlackChannelDoNotExistError
	ExitCodeFailedToFetchGitHubStatusError
	ExitCodeFailedToPostToSlackError
)

// PrintGreenF prints green success message on console
func PrintGreenF(writer io.Writer, format string, args ...interface{}) {
	PrintColorF(writer, "green", format, args...)
}

// PrintRedF prints red error message on console
func PrintRedF(writer io.Writer, format string, args ...interface{}) {
	PrintColorF(writer, "red", format, args...)
}

// PrintColorF prints colored message on console
func PrintColorF(writer io.Writer, color, format string, args ...interface{}) {
	format = fmt.Sprintf("[%s]%s[reset]", color, format)
	fmt.Fprint(writer, colorstring.Color(fmt.Sprintf(format, args...)))
}

// DebugF prints debug message when EnvDebug is true
func DebugF(format string, args ...interface{}) {
	if env := os.Getenv(EnvDebug); len(env) != 0 {
		log.Printf("[DEBUG] "+format+"\n", args...)
	}
}

func getenvOrDefault(key string, defvalue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	} else {
		return defvalue
	}
}

func getenvOrExit(key string, code int) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Println(fmt.Sprintf("Fail to look up %s.", key))
		os.Exit(code)
	}
	return value
}

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *strslice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		version bool

		users             strslice
		channel           string
		webhookURL        string
		name              string
		iconEmojiPositive string
		iconEmojiNegative string
		msgNegative       string
	)

	flags := flag.NewFlagSet(Name, flag.ExitOnError)
	flags.SetOutput(cli.errStream)

	flags.Var(&users, "user", "GitHub user names")
	flags.Var(&users, "u", "GitHub user names")

	flags.StringVar(&channel, "channel", "", "Target Slack channel")
	flags.StringVar(&channel, "c", "", "Target Slack channel")

	flags.StringVar(&webhookURL, "webhook-url", "", "Slack incoming Webhook URL")
	flags.StringVar(&webhookURL, "w", "", "Slack incoming Webhook URL")

	flags.StringVar(&name, "name", "kusabot", "Bot name")
	flags.StringVar(&name, "n", "kusabot", "Bot name")

	flags.StringVar(&iconEmojiPositive, "icon-emoji-positive", ":seedling:", "")
	flags.StringVar(&iconEmojiNegative, "icon-emoji-negative", ":japanese_goblin:", "")
	flags.StringVar(&msgNegative, "message-negative", ":warning: There are no contributions today ! :warning:", "")

	flags.BoolVar(&version, "version", false, "")
	flags.BoolVar(&version, "v", false, "")

	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeParseFlagsError
	}

	if version {
		fmt.Fprintf(cli.outStream, OutputVersion())
		return ExitCodeOK
	}

	if len(users) < 1 {
		PrintRedF(cli.errStream, "You must set 1 GitHub user at least with `--user` option.")
		os.Exit(ExitCodeGitHubUsersDoNotExistError)
	}

	if len(webhookURL) < 1 {
		fmt.Print("Incoming Webhook URL: ")
		fmt.Scan(&webhookURL)
		if len(webhookURL) < 1 {
			PrintRedF(cli.errStream, "You must set slack webhook url with `--url` option or tty.")
			os.Exit(ExitCodeSlackWebhookUrlDoNotExistError)
		}
	}

	kusa := &Kusa{}
	statuses, err := kusa.Fetch(users)

	if err != nil {
		PrintRedF(cli.errStream, "%v", err)
		return ExitCodeFailedToFetchGitHubStatusError
	}

	slack := &Slack{
		WebhookUrl:               webhookURL,
		Channel:                  channel,
		User:                     name,
		IconEmoji:                iconEmojiPositive,
		IconEmojiNoContributions: iconEmojiNegative,
		MsgNoContributions:       msgNegative,
	}

	for _, status := range statuses {
		err := slack.Post(status)
		if err != nil {
			PrintRedF(cli.errStream, "%v", err)
			return ExitCodeFailedToPostToSlackError
		}
	}

	return ExitCodeOK
}
