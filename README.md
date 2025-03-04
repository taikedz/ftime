# `ftime`

A simple program to retain the walltime of a command.

## Usage

```sh
# Structure
ftime [-t FILE] -- COMMAND ...


# Example
# store timings by default in 't.times'
ftime -- find ~/ -name '*.mp3'

# store timings in 'durations.txt'
ftime -t durations.txt -- find ~/ -name '*.mp3'
```

## Motivation

The Linux/UNIX `time` command sometimes behaves erratically in its output.

It should be possible to achieve similar effect with `time -o FILE --portability --append COMMAND ...`

I just needed an exercise example whilst learning go.
