package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(mlCmd)
	mlCmd.Flags().StringP("json-file", "j", "", "file with list of packages (json format)")
	if err := viper.BindPFlag("json-file", mlCmd.Flags().Lookup("json-file")); err != nil {
		panic(err)
	}
}

var mlCmd = &cobra.Command{
	Use:   "ml",
	Short: "ml",
	RunE:  ml,
}

func ml(cmd *cobra.Command, args []string) error {
	type pkg struct {
		Name  string
		Owner string
		GoSig string
	}
	var packages []pkg

	file, _ := ioutil.ReadFile(viper.GetString("json-file"))

	_ = json.Unmarshal([]byte(file), &packages)

	var pkgs []pkg
	for _, p := range packages {
		if p.GoSig != "admin" && p.GoSig != "commit" {
			pkgs = append(pkgs, p)
		}
	}

	fmt.Printf("================================================================================\n\n")
	fmt.Printf("Maintainers per package:\n\n")

	for i := 0; i < len(pkgs); i++ {
		fmt.Printf(" - %v: %v\n", pkgs[i].Name, pkgs[i].Owner)
	}
	fmt.Printf("\n\nPackages per maintainer:\n\n")

	packagers := make(map[string][]string)
	for i := 0; i < len(pkgs); i++ {
		v, e := packagers[pkgs[i].Owner]
		if e {
			v = append(v, pkgs[i].Name)
			packagers[pkgs[i].Owner] = v
		} else {
			packagers[pkgs[i].Owner] = []string{pkgs[i].Name}
		}
	}

	keys := make([]string, 0, len(packagers))
	for k := range packagers {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%v (%v): %v\n\n", k, len(packagers[k]), strings.Join(packagers[k], ", "))
	}

	return nil
}
