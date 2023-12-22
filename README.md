# Mediocre Word Count (mwc)
A simple [word count](https://en.wikipedia.org/wiki/Wc_(Unix)) clone in golang.
This is my solution of [this coding challenge](https://codingchallenges.fyi/challenges/challenge-wc/). 
Feedback and stuff is appreciated 😘

# Running
Make sure you have [golang](https://go.dev/) installed
```sh
# build the program
go build -o bin/mwc main.go
```
```sh
# running
# ./bin/mwc [OPTIONS] [FILE]
./bin/mwc ./test.txt    # prints the line, word and byte count
./bin/mwc -c ./test.txt # prints the byte count
./bin/mwc -l ./test.txt # prints the line count
./bin/mwc -L ./test.txt # print the length of the longest line
./bin/mwc -m ./test.txt # prints the character count
./bin/mwc -w ./test.txt # prints the word count

# can also read input from stdin
cat ./test.txt | ./bin/mvc # equivalent to `./bin/mwc ./test.txt`
```

## Todo
- [x] '-L' option (print the length of the longest line)
- [ ] help option (maybe use flag or something for reading args)
- [ ] that total count thing wc does when you give it multiple files
