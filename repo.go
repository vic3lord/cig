package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func checkRepo(root string, path string, channel chan string, wg *sync.WaitGroup) {
	exists, err := exists(fmt.Sprintf("%s%s.git", path, string(os.PathSeparator)))

	if exists {
		modified_files := exec.Command("git", "status", "-s")
		modified_files.Dir = path

		count_out, _ := modified_files.Output()
		modified_lines := strings.Split(string(count_out), "\n")
		modified := len(modified_lines) - 1

		if err != nil {
			println(err.Error())
			return
		}

		changes := []string{}

		if modified > 0 && modified_lines[0] != "" {
			changes = append(changes, print_output(fmt.Sprintf(" M(%d)", modified), "red"))
		}

		branch := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branch.Dir = path
		bstdout, _ := branch.Output()
		branch_name := strings.TrimSpace(string(bstdout[:]))

		local := exec.Command("git", "rev-parse", branch_name)
		local.Dir = path
		lstdout, _ := local.Output()
		local_ref := strings.TrimSpace(string(lstdout[:]))

		remote := exec.Command("git", "rev-parse", fmt.Sprintf("origin/%s", branch_name))
		remote.Dir = path
		rstdout, err := remote.Output()
		remote_ref := strings.TrimSpace(string(rstdout[:]))

		if err == nil && remote_ref != local_ref {
			changes = append(changes, print_output(" P", "blue"))
		}

		if len(changes) > 0 {
			var buffer bytes.Buffer

			repo_name := strings.Replace(path, fmt.Sprintf("%s%s", root, string(os.PathSeparator)), "", -1)

			buffer.WriteString(fmt.Sprintf("- %s (%s)", repo_name, branch_name))
			for _, change := range changes {
				buffer.WriteString(change)
			}
			channel <- buffer.String() + "\n"
		}

	}
	wg.Done()
}
