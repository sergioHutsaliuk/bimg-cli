/*
Copyright © 2023 Serhii Hutsaliuk <gutsalyuk.sergio@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// compressCmd represents the compress command
var compressCmd = &cobra.Command{
	Use:   "compress",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(1)),
	RunE: func(cmd *cobra.Command, args []string) error {
		return compress(args[0])
	},
}

func init() {
	rootCmd.AddCommand(compressCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compressCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// compressCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func compress(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	err = os.Mkdir("compressed", 0750)
	if !errors.Is(err, os.ErrExist) {
		return err
	}

	for _, file := range files {
		filename := file.Name()

		if strings.HasPrefix(filename, ".") {
			continue
		}

		buf, bufErr := bimg.Read(fmt.Sprintf("%s/%s", dir, filename))
		if bufErr != nil {
			return bufErr
		}

		processed, procErr := bimg.NewImage(buf).Process(bimg.Options{Quality: 40})
		if procErr != nil {
			return procErr
		}

		if errWrite := bimg.Write(fmt.Sprintf("./compressed/%s", file.Name()), processed); errWrite != nil {
			log.Printf("Error %v in processing %s\n", errWrite, filename)

			continue
		}

		log.Printf("Successful - %s\n", filename)
	}

	log.Printf("All processed\n")

	return nil
}
