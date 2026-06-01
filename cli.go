package main

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/urfave/cli/v3"
)

// Primary CLI command for whoistrader, taking in multiple steamIDs. Requires a pre-made registry of endpoints.
func buildCLI(registry *Registry) *cli.Command {
	return &cli.Command{
		Name:  "whoistrader",
		Usage: "CS2 trader profiler to vet players",
		Flags: []cli.Flag{
			&cli.StringFlag {
				Name: "output",
				Aliases: []string{"o"},
				Usage: "Write output to `filepath`. If it does not exist, it will be automatically created.",
			},
		},
		Commands: []*cli.Command{
			{
				Name:      "profile",
				Usage:     "Fetch aggregated trader profile data from one or more Steam IDs",
				ArgsUsage: "<steamID> [steamID...] --output [filepath]",
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if cmd.Args().Len() == 0 {
						return fmt.Errorf("Command requires at least one Steam ID")
					}

					for _, arg := range cmd.Args().Slice() {
						//uint64 SteamID
						steamID, err := strconv.ParseUint(arg, 10, 64)
						if err != nil {
							return fmt.Errorf("Invalid Steam ID: %s", arg)
						}

						if err := CreateProfile(ctx, steamID, registry, cmd.Root().String("output")); err != nil {
							log.Printf("Failed to create profile for steamID %d: %s\n", steamID, err)
						}
					}
					return nil
				},
			},
		},
	}
}
