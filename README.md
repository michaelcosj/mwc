# Mediocre Word Count (mwc)
A simple [word count](https://en.wikipedia.org/wiki/Wc_(Unix)) in golang.
This is a solution of [this coding challenge](https://codingchallenges.fyi/challenges/challenge-wc/). 
Feedback and stuff is appreciated 😘

# Running
Make sure you have [golang](https://go.dev/) installed

```sh
    go build -o bin/mwc main.go
```
```sh
    ./bin/mwc [OPTIONS] [FILE]
```

Options are the basic ones from wc. It can also read from stdin

## Todo
- [] help option (maybe use clap or something for reading args)
- [] that total count thing wc does when you give it multiple files
- [] '-L' option from wc
