// Copyright © 2018 Eric Sizer <eric.sizer@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

var (
	name      string
	site      string
	directory string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Generate a new WordPress install",
	Long: `Download the latest WordPress and install the FP theme.
	Once complete, cd into the theme and run npm install.
	`,
	Run: install,
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVarP(&name, "name", "n", "fp", "Name of theme (required)")
	newCmd.Flags().StringVarP(&site, "site", "s", "", "Site name (WordPress Folder")
	newCmd.Flags().StringVarP(&directory, "dir", "d", "./", "Install to")

	newCmd.MarkFlagRequired("name")
}

// Install WordPress
func install(ccmd *cobra.Command, args []string) {
	if site == "" {
		site = name
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	defer returnToCwd(cwd)

	getWordPress()
}

// Get the latest wordpress then unzip it
func getWordPress() {

	defer cleanUp()

	// Download the latest WordPress
	wget := exec.Command("wget", "http://wordpress.org/latest.tar.gz")
	runCmd(wget, "Downloading Wordpress")

	// Unzip the tar.gz
	tar := exec.Command("tar", "xfz", "latest.tar.gz")
	runCmd(tar, "Unzipping")

	// Move to specific directory
	target := path.Join(directory, site)
	empty, _ := isEmpty(target)
	if empty {
		mv := exec.Command("mv", "./wordpress")
		runCmd(mv, fmt.Sprintln("Moving Wordpress to ", path.Join(directory, site)))
	} else {
		log.Println("Target directory is not empty.")
	}

}

func runCmd(cmd *exec.Cmd, message string) {
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	if message != "" {
		log.Println(message)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func cleanUp() {
	// Clean up after yourself
	cleanUprmDir := exec.Command("rm", "-rf", "./wordpress")
	runCmd(cleanUprmDir, "Cleaning up files")

	cleanUprm := exec.Command("rm", "-f", "latest.tar.gz")
	runCmd(cleanUprm, "")
}

func returnToCwd(cwd string) {

	cmd := exec.Command("cd", cwd)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

}

func isEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err // Either not empty or error, suits both cases
}
