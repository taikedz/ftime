/* Print the wall-clock time of the duration of a subprocess.
 *
 * Linux `time` command frequently writes strange output formatting on console if attempting to store
 * the result, and it's sometimes also hard to dissociate it from the program outputs.
 *
 * This simple program runs a command, continually prints its output, and then writes the duration to a file.
 */

package main

import (
    "fmt"
    "time"
    "os/exec"
    "os"
    "bufio"
    "log"
    "strings"
    "io"
)

func main() {
    if len(os.Args) < 3 {
        fmt.Printf("%s TIMERFILE COMMAND ...\n", os.Args[0])
        os.Exit(1)
    }

    t0 := time.Now()

    status := runCommand(os.Args[2], os.Args[3:]...)

    t1 := time.Now()
    writeTime(os.Args[1], t1.Sub(t0), strings.Join(os.Args[2:], " "))

    os.Exit(status)
}

func writeTime(filename string, tdiff time.Duration, command string) {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        log.Fatal("Could not open timings.txt file")
    }
    defer f.Close()

    message := fmt.Sprintf("[%v] : %s\n", tdiff, command)

    if _, err = f.WriteString(message); err != nil {
        log.Fatal("Could not write timing")
    }
}

/* Simulat a non-capturing shell command execution
 */
func runCommand(command string, args... string) int {
    cmd := exec.Command(command, args...)
    stdout, err := cmd.StdoutPipe()
    if err != nil {
        log.Fatal(err)
    }
    cmd.Start()

    reader := bufio.NewReader(stdout)
    buffer := make([]byte, 1024)
    for {
        /* FIXME - this Read(buf) will attempt to wait for a full buffer's worth before returning
         * Ideally, should return whatever's currently available (e.g. if subprocessed is momentarily not outputting)
         */
        count, readerr := reader.Read(buffer)
        if readerr != nil && readerr != io.EOF {
            fmt.Printf("! %v\n", readerr)
            log.Fatal("Process output reading failed")
        }

        if count == 0 {
            if err := cmd.Wait(); err != nil {
                // Linux-only. Unsure what this does on Windows
                if exiterr, ok := err.(*exec.ExitError); ok {
                    return exiterr.ExitCode()
                } else {
                    return 1
                }
            }
            return 0
        } else {
            fmt.Printf("%s", buffer[:count])
        }
    }
}
