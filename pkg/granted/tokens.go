package granted

import (
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-fate/granted/pkg/credstore"
	"github.com/common-fate/granted/pkg/testable"
	"github.com/urfave/cli/v2"
)

var TokenCommand = cli.Command{
	Name:        "token",
	Usage:       "Manage aws access tokens",
	Subcommands: []*cli.Command{&TokenListCommand, &ClearTokensCommand, &ClearAllTokensCommand},
	Action: func(c *cli.Context) error {

		return nil
	},
}

var TokenListCommand = cli.Command{
	Name:  "list",
	Usage: "Lists all access tokens saved in the keyring",
	Action: func(ctx *cli.Context) error {
		tokens, err := credstore.List()
		if err != nil {
			return err
		}
		for i, token := range tokens {
			fmt.Fprintf(os.Stderr, "%d. %s---(%s)\n", i+1, token.Key, token.Description)
		}
		return nil
	},
}
var ClearAllTokensCommand = cli.Command{
	Name:  "clear",
	Usage: "Remove all saved tokens from keyring",
	Action: func(c *cli.Context) error {
		err := credstore.ClearAll()
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Cleared all saved tokens")
		return nil
	},
}

var ClearTokensCommand = cli.Command{
	Name:  "remove",
	Usage: "Remove a selected token from the keyring",
	Action: func(c *cli.Context) error {
		var selection string

		if c.Args().First() != "" {
			selection = c.Args().First()
		}

		keys, err := credstore.List()
		if err != nil {
			return err
		}
		if selection == "" {
			tokenList := []string{}
			for _, t := range keys {
				tokenList = append(tokenList, t.Key+"---"+t.Description+")")
			}
			withStdio := survey.WithStdio(os.Stdin, os.Stderr, os.Stderr)
			in := survey.Select{
				Message: "Select a token to remove from keyring",
				Options: tokenList,
			}
			fmt.Fprintln(os.Stderr)
			err = testable.AskOne(&in, &selection, withStdio)
			if err != nil {
				return err
			}
		}
		selection = strings.Split(selection, "---")[0]

		err = credstore.Clear(selection)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "Cleared %s", selection)
		return nil
	},
}
