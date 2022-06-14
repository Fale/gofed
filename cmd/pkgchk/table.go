package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fale/gofed/pkg/oraculum"
	"github.com/fale/gofed/pkg/pagure"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(tableCmd)
	tableCmd.Flags().StringP("packages-file", "f", "", "file with list of packages (one per line)")
	if err := viper.BindPFlag("packages-file", tableCmd.Flags().Lookup("packages-file")); err != nil {
		panic(err)
	}
}

var tableCmd = &cobra.Command{
	Use:   "table",
	Short: "table",
	RunE:  table,
}

func table(cmd *cobra.Command, args []string) error {
	file, err := os.Open(viper.GetString("packages-file"))
	if err != nil {
		return err
	}
	defer file.Close()

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package", "Owner", "GoSig", "FTBFS"})
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pars := pagure.GetProjectsParameters{
			BaseUrl:   "https://src.fedoraproject.org",
			NameSpace: "rpms",
			Pattern:   scanner.Text(),
			Fork:      "false",
		}
		projects, err := pagure.GetProjects(pars)
		if err != nil {
			return err
		}
		for _, project := range projects {
			fmt.Printf("%v\n", project.Name)
			var gslvl string
			for g, v := range project.AccessGroups {
				for _, n := range v {
					if n == "go-sig" {
						gslvl = g
					}
				}
			}
			p, err := oraculum.LookupPackage(project.Name)
			if err != nil {
				return err
			}
			var ftbfs37 string
			if p.BranchFTBFS("Fedora Rawhide") {
				ftbfs37 = "YES"
			}
			table.Append([]string{project.Name, project.User.Name, gslvl, ftbfs37})
		}
	}
	table.Render()

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
