package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/fale/gofed/pkg/pagure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(jsnCmd)
	jsnCmd.Flags().StringP("package-file", "f", "", "file with list of packages (one per line)")
	jsnCmd.Flags().StringP("output-file", "o", "", "output file")
	if err := viper.BindPFlag("package-file", jsnCmd.Flags().Lookup("package-file")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("output-file", jsnCmd.Flags().Lookup("output-file")); err != nil {
		panic(err)
	}
}

var jsnCmd = &cobra.Command{
	Use:   "json",
	Short: "json",
	RunE:  jsn,
}

func jsn(cmd *cobra.Command, args []string) error {
	type pkg struct {
		Name  string
		Owner string
		GoSig string
	}
	var pkgs []pkg
	fmt.Println(viper.GetString("package-file"))
	fmt.Println(viper.GetString("output-file"))
	file, err := os.Open(viper.GetString("package-file"))
	if err != nil {
		return err
	}
	defer file.Close()

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
			fmt.Println(project.Name)
			var gslvl string
			for g, v := range project.AccessGroups {
				for _, n := range v {
					if n == "go-sig" {
						gslvl = g
					}
				}
			}
			p := pkg{
				Name:  project.Name,
				Owner: project.User.Name,
				GoSig: gslvl,
			}
			pkgs = append(pkgs, p)
		}
	}

	f, _ := json.MarshalIndent(pkgs, "", " ")

	_ = ioutil.WriteFile(viper.GetString("output-file"), f, 0644)

	return nil
}
