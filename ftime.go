/* Print the wall-clock time of the duration of a subprocess.
 *
 * Linux `time` command frequently writes strange output formatting on console if attempting to store
 * the result, and it's sometimes also hard to dissociate it from the program outputs.
 *
 * This simple program runs a command, dumping output, and then writes the duration to a specified file.
 */

package main

import (
    "fmt"
    "time"
    "os/exec"
    "os"
    "log"
    "strings"
    "path/filepath"
)

const ERR_UNKNOWN = 120

func main() {
    if len(os.Args) < 3 {
        binname := filepath.Base(os.Args[0])
        fmt.Printf("%s TIMERFILE COMMAND ...\n", binname)
        fmt.Printf("e.g.\n")
        fmt.Printf("    %s t.times find ~/ -name '*.mp3'\n", binname)
        os.Exit(1)
    }

    // Pre-emptively create file to be sure we can actually access it
    //  (at least at ths start...)
    f := openFile(os.Args[1])
    f.Close()

    t0 := time.Now()
    status := runCommand(os.Args[2], os.Args[3:]...)
    t1 := time.Now()

    writeTime(os.Args[1], t1.Sub(t0), strings.Join(os.Args[2:], " "))

    os.Exit(status)
}

func openFile(filename string) *os.File {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        log.Fatal(fmt.Sprintf("Could not open '%s' file", filename))
    }

    return f
}

func writeTime(filename string, tdiff time.Duration, command string) {
    f := openFile(filename)
    defer f.Close()

    message := fmt.Sprintf("[%v] : %s\n", tdiff, command)

    if _, err := f.WriteString(message); err != nil {
        log.Fatal("Could not write timing")
    }
}

func runCommand(command string, args... string) int {
    cmd := exec.Command(command, args...)

    // https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr

    // We need to use Start() if we want to later use Wait() and get the exit status information
    //   else we could have just used Run()
    err := cmd.Start()
    if err != nil {
        log.Fatal(err)
    }

    // https://stackoverflow.com/a/10385867/2703818
    if err = cmd.Wait(); err != nil {
        if exiterr, ok := err.(*exec.ExitError); ok {
            return exiterr.ExitCode()
        } else {
            return ERR_UNKNOWN
        }
    }
    return 0
}

