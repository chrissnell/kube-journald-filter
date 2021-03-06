package main

// #cgo pkg-config: ---cflags --libs libsystemd
// #include <systemd/sd-journal.h>

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/coreos/go-systemd/sdjournal"
	"github.com/coreos/go-systemd/util"
	"k8s.io/kubernetes/pkg/kubelet/dockertools"
)

func main() {
	f := flag.NewFlagSet("kube-journald-filter", flag.ExitOnError)
	journalPath := f.String("alt-journal-base", "", "Use alternate base directory for journal.  Directory will be appended with /<machine-id> automatically.")
	f.Parse(os.Args[1:])

	var j *sdjournal.Journal
	var err error

	if len(*journalPath) > 0 {
		// Open a journal in a custom location

		// Get our machine-id, used to choose the correct journal subdirectory
		machineID, err := util.GetMachineID()
		if err != nil {
			log.Fatalln("Could not get machine-id:", err)
		}
		fullPath := filepath.Join(*journalPath, machineID)

		// Open the journal in our alternate location
		j, err = sdjournal.NewJournalFromDir(fullPath)
		if err != nil {
			log.Fatalln("Could not open journal:", err)
		}
	} else {
		// Open the system journal in the standard location
		j, err = sdjournal.NewJournal()
		if err != nil {
			log.Fatalln("Could not open journal:", err)
		}
	}

	// Seek to the end of the journal and begin "tailing" from there
	err = j.SeekTail()
	if err != nil {
		log.Fatalln("Could not seek to tail of journal:", err)
	}

	// Wait for new entries to show in the journal
	for {
		n, err := j.Next()
		if n < 1 {
			j.Wait(sdjournal.IndefiniteWait)
			continue
		}
		if err != nil {
			log.Println("Could not advance next read pointer in journald:", err)
			continue
		}

		// A journal entry is really just a set of key-values.  The actual content of the logged message
		// is stored under the MESSAGE key.  Retrieve it.
		msg, err := j.GetDataValue("MESSAGE")
		if err != nil {
			log.Println("Could not read MESSAGE from journald", err)
			continue
		}

		hostname, err := j.GetDataValue("_HOSTNAME")
		if err != nil {
			log.Println("Could not read _HOSTNAME from journald", err)
			continue
		}

		pid, err := j.GetDataValue("_PID")
		if err != nil {
			log.Println("Could not read _PID from journald", err)
			continue
		}

		cmd, err := j.GetDataValue("_COMM")
		if err != nil {
			log.Println("Could not read _COMM from journald", err)
			continue
		}

		// Docker logs the container name to CONTAINER_NAME.  Fetch that and if it's present, discard the prefix.
		containerName, _ := j.GetDataValue("CONTAINER_NAME")
		if len(containerName) > 0 {
			// Kubernetes-launched containers always start with "k8s".  This is a constant specified in the
			// Kubernetes source code.  Unfortunately, it's not exported from there so we can't draw upon it
			// and must define it here.
			if strings.HasPrefix(containerName, "k8s") {
				// We use Kubernetes' library to parse the name.  It's really just a simple regex but the
				// format might change in the future so this is the easiest way to retain compatability.
				kcn, _, err := dockertools.ParseDockerName(containerName)
				if err != nil {
					log.Println("Error parsing container name.")
				} else {
					tr := regexp.MustCompile(`(.*)_([a-z0-9-]*)$`)
					matches := tr.FindStringSubmatch(kcn.PodFullName)
					fmt.Printf("%v %v[%v] NS=%v POD=%v %v\n", hostname, cmd, pid, matches[2], matches[1], msg)
				}
			}
		} else {
			// If we didn't see a CONTAINER_NAME, just print the message as-is
			fmt.Printf("%v %v[%v] %v\n", hostname, cmd, pid, msg)
		}
	}

}
