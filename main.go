package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Command interface {
	Exec(handler Handler)
}
type Handler interface {
	Post(cmd Command)
}

type PrintCommand struct {
	arg string
}

func (cmd *PrintCommand) Exec(loop Handler) {
	fmt.Println(cmd.arg)
}

type AddCommand struct {
	arg1, arg2 int
}

func (cmd *AddCommand) Exec(loop Handler) {
	loop.Post(&PrintCommand{arg: strconv.Itoa(cmd.arg1 + cmd.arg2)})
}

type EventLoop struct {
	queue []Command
}

func (eventLoop *EventLoop) Post(cmd Command) {
	eventLoop.queue = append(eventLoop.queue, cmd)
}

func (eventLoop *EventLoop) Start() {
	eventLoop.queue = make([]Command, 0)
}

func (eventLoop *EventLoop) AwaitFinish() {
	for len(eventLoop.queue) > 0 {
		cmd := eventLoop.queue[0]
		eventLoop.queue = eventLoop.queue[1:]
		cmd.Exec(eventLoop)
	}
}

func parse(line string) Command {
	parts := strings.Fields(line)
	cmd := parts[0]
	args := parts[1:]
	switch cmd {
	case "print":
		return &PrintCommand{arg: strings.Join(args, " ")}
	case "add":
		arg1, err := strconv.Atoi(args[0])
		if err != nil {
			arg := fmt.Sprintf("Syntax Error: first argument is not an integer, %s", err)
			return &PrintCommand{arg}
		}

		arg2, err := strconv.Atoi(args[1])
		if err != nil {
			arg := fmt.Sprintf("Syntax Error: second argument is not an integer, %s", err)
			return &PrintCommand{arg}
		}

		if len(args) > 2 {
			arg := fmt.Sprintf("Syntax Error: 'add' command accepts only 2 arguments")
			return &PrintCommand{arg}
		}

		return &AddCommand{arg1, arg2}
	}
	return nil
}

func main() {
	eventLoop := new(EventLoop)
	inputFile := "input.txt"

	eventLoop.Start()
	fmt.Println("Event loop started")
	fmt.Println()

	if input, err := os.Open(inputFile); err == nil {
		defer input.Close()
		scanner := bufio.NewScanner(input)

		for scanner.Scan() {
			commandLine := scanner.Text()
			cmd := parse(commandLine) // parse the line to get a Command

			if cmd != nil {
				eventLoop.Post(cmd)
			}
		}
	}

	eventLoop.AwaitFinish()

	fmt.Println()
	fmt.Println("Event loop finished")
}
