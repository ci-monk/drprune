package github

import (
	"bufio"
	"context"
	"fmt"
	"os"

	tm "github.com/buger/goterm"
	"github.com/google/go-github/v41/github"
	"github.com/lpmatos/drprune/internal/log"
	"github.com/lpmatos/drprune/internal/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func NewCmdInsights() *cobra.Command {
	var insightsCmd = &cobra.Command{
		Use:   "insights",
		Short: "Get insights of GitHub Registry (ghcr.io)",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			container = utils.EncodeParam(container)
			totalPackages := 0

			// Auth in github client
			ctx := context.Background()
			ts := oauth2.StaticTokenSource(
				&oauth2.Token{AccessToken: token},
			)
			tc := oauth2.NewClient(ctx, ts)
			client := github.NewClient(tc)

			// Get all packages of user.
			pkgs, _, err := client.Users.ListPackages(ctx, name, &github.PackageListOptions{
				PackageType: github.String("container"),
				ListOptions: github.ListOptions{
					PerPage: 300,
				},
			})
			if err != nil {
				log.Errorf("List package err: %v", err)
			}

			// Check if user have packages.
			if len(pkgs) > 0 {
				// If user have packages, loop.
				for _, pkg := range pkgs {
					totalTags, totalUntagged := 0, 0
					fmt.Printf("\n\n=================================================\n\n")
					pterm.DefaultSection.Println("Package Information")
					fmt.Printf("> Package name: %v\n", pkg.GetName())
					fmt.Printf("> Package id: %d\n", pkg.GetID())
					fmt.Printf("> Package owner: %s\n", pkg.GetOwner().GetLogin())

					// Get all versions of the package.
					pkgVersions, _, err := client.Users.PackageGetAllVersions(ctx,
						name, "container", utils.EncodeParam(pkg.GetName()),
					)
					if err != nil {
						log.Fatal(err)
					}

					pterm.Println()
					pterm.DefaultSection.Println("Package Versions Information")

					// Loop each version of the package.
					for _, pkgVersion := range pkgVersions {
						fmt.Printf("\n> Package version name: %s\n", *pkgVersion.Name)
						fmt.Printf("> Package version id: %d\n", *pkgVersion.ID)

						tags := pkgVersion.GetMetadata().GetContainer().Tags
						if len(tags) == 0 {
							totalUntagged++
							fmt.Printf("> Empty tags\n")
							continue
						}

						fmt.Printf("> Package tags: %v\n", tags)
						for i := 0; i < len(tags); i++ {
							totalTags++
						}
					}

					pterm.Println()
					pterm.DefaultSection.Println("Resume Information")
					fmt.Printf("----> Total tagged for %s package: %v", pkg.GetName(), totalTags)
					fmt.Printf("\n----> Total untagged for %s package: %v", pkg.GetName(), totalUntagged)
					totalPackages++

					consoleReader := bufio.NewReaderSize(os.Stdin, 1)
					fmt.Print("\n\n> Press [ENTER] to see the next or [ESC] to exit: ")
					input, _ := consoleReader.ReadByte()
					ascii := input
					if ascii == 27 {
						fmt.Printf("\nExiting...\n")
						os.Exit(0)
					}

					tm.Clear()
					tm.MoveCursor(1, 1)
					tm.Flush()
				}
			} else {
				log.Warnf("The user %s don't any tags!", name)
			}

			fmt.Printf("\n\n=================================================\n")
			pterm.Println()
			pterm.DefaultSection.Println("Final Information")
			fmt.Printf("\n----> The user %s have %v packages\n", name, totalPackages)
		},
	}
	return insightsCmd
}
