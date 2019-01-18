package main

import (
	"fmt"
	"os"
	"strings"
	"time"
	"bufio"
	// "io"
	"io/ioutil"
	"log"

	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

func check(e error) {
	if e != nil {
			panic(e)
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func readFile(filePath string, shell *ishell.Shell) {
	data, err := ioutil.ReadFile(filePath)
	check(err)
	shell.Println(string(data))
}


func writeFile(filePath string, content string) {
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
			log.Fatal(err)
	}
	if _, err := f.Write([]byte(content + "\n")); err != nil {
			log.Fatal(err)
	}
	if err := f.Close(); err != nil {
			log.Fatal(err)
	}

}

func returnJobChoices(filePath string, shell *ishell.Shell) []string {
	file, err := os.Open(filePath)
	if err != nil {
			log.Fatal(err)
	}
	defer file.Close()
	
	var jobs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
			item := scanner.Text()
			formattedItem := strings.TrimSpace(strings.Split(item, ",")[2])
			jobs = append(jobs, formattedItem)
	}

	if err := scanner.Err(); err != nil {
			log.Fatal(err)
	}

	return jobs
}


func clockIn(job string, shell *ishell.Shell) {
	currentJobTimeFile, err := os.Open("./text_files/current_job_time.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer currentJobTimeFile.Close()

	fileStat, statErr := currentJobTimeFile.Stat()
	if statErr != nil {
		log.Fatal(statErr)
		// Could not obtain stat, handle error
	}

	if fileStat.Size() > 0 {
			// clock out of current job and add it to rolling timesheet
			// delete that line from the file
	}


}

func main() {
	// TODO: Command to display stats
	shell := ishell.New()

  t := time.Now()
	currentTime := t.Format("2006-01-02 15:04:05")
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	boldRed := color.New(color.FgRed, color.Bold).SprintFunc()

	// display info.
	shell.Println("--------->", cyan("Welcome To Time Muncher")," <----------")
	shell.Println(" ")
	shell.Println("It's", currentTime)
	shell.Println(" ")
	shell.Println("You can begin by asking for", yellow("help"))
	shell.Println(" ")

	// handle "greet".
	shell.AddCmd(&ishell.Cmd{
		Name:    "add-job",
		Help:    `Add a new job! Give it a formal name, and the quickbook id, short name, `,
		Func: func(c *ishell.Context) {
			name := ""
			if len(c.Args) >= 3 {
				name = strings.Join(c.Args, " ")
				writeFile("./text_files/jobs_and_codes.txt", name)
				c.Println("You've added", name)
			} else {
				c.Println("Ahh.. It looks like you didn't provide enough information.")
				c.Println("Be sure to include the formal name, the quickbook id, and the short name")
				c.Println("for example: Critique LLC, 40-R3-Simon:41-R3-Josh, Round3")
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "read",
		Help: "Read the contents of the rolling-timesheet",
		Func: func(c *ishell.Context) {
			file := ""

			if len(c.Args) == 1 {
				file = strings.Join(c.Args, " ")
			}
			switch file {
			case "rolling_timesheet":
				readFile("./text_files/rolling_timesheet.txt", shell)
				break;
			case "jobs_and_codes":
				readFile("./text_files/jobs_and_codes.txt", shell)
				break;
			case "user_info":
				readFile("./text_files/user_info.txt", shell)
				break;
			case "example_timesheet":
				readFile("./text_files/example_timesheet.iif", shell)
				break;
			default:
				shell.Print("No file with that name was found")
			}
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "time",
		Help: "Returns the current time",
		Func: func(c *ishell.Context) {
			t := time.Now()
			currentTime := t.Format("15:04:05")
			shell.Println("The Current Time is", currentTime)
			
		},
	})

	// read multiple lines with "multi" command
	shell.AddCmd(&ishell.Cmd{
		Name: "multi",
		Help: "input in multiple lines",
		Func: func(c *ishell.Context) {
			c.Println("Input multiple lines and end with semicolon ';'.")
			lines := c.ReadMultiLines(";")
			c.Println("Done reading. You wrote:")
			c.Println(lines)
		},
	})

	// multiple choice
	shell.AddCmd(&ishell.Cmd{
		Name: "clock-in",
		Help: "multiple choice prompt",
		Func: func(c *ishell.Context) {
			jobs := returnJobChoices("./text_files/jobs_and_codes.txt", shell)
			choice := c.MultiChoice(jobs, "Who would you like to clock into?")
			clockIn(jobs[choice], shell)
			c.Println("You're clocked into", jobs[choice])
		},
	})

	// multiple choice
	shell.AddCmd(&ishell.Cmd{
		Name: "in",
		Help: "multiple choice prompt (short for clock-in)",
		Func: func(c *ishell.Context) {
			jobs := returnJobChoices("./text_files/jobs_and_codes.txt", shell)
			choice := c.MultiChoice(jobs, "Who would you like to clock into?")
			clockIn(jobs[choice], shell)
			c.Println("You're clocked into", jobs[choice])
		},
	})

	// multiple choice
	shell.AddCmd(&ishell.Cmd{
		Name: "edit-job",
		Help: "checklist prompt",
		Func: func(c *ishell.Context) {
			languages := []string{"Radius", "Round3", "Bee Corp", "Moana"}
			choices := c.Checklist(languages,
				"Which job would you like to edit?",
				nil)
			out := func() (c []string) {
				for _, v := range choices {
					c = append(c, languages[v])
				}
				return
			}
			c.Println("You've selected'", strings.Join(out(), ", "))
		},
	})

	// progress bars
	{

		// indeterminate
		shell.AddCmd(&ishell.Cmd{
			Name: "create-timesheet",
			Help: "Creates an .iif Timesheet",
			Func: func(c *ishell.Context) {
				c.ProgressBar().Indeterminate(true)
				c.ProgressBar().Start()
				time.Sleep(time.Second * 10)
				c.ProgressBar().Stop()
			},
		})
	}

	shell.AddCmd(&ishell.Cmd{
		Name: "paged",
		Help: "show paged text",
		Func: func(c *ishell.Context) {
			lines := ""
			line := `%d. This is a paged text input.
This is another line of it.
`
			for i := 0; i < 100; i++ {
				lines += fmt.Sprintf(line, i+1)
			}
			c.ShowPaged(lines)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "color",
		Help: "color print",
		Func: func(c *ishell.Context) {
			c.Print(cyan("cyan\n"))
			c.Println(yellow("yellow"))
			c.Printf("%s\n", boldRed("bold red"))
		},
	})

	// when started with "exit" as first argument, assume non-interactive execution
	if len(os.Args) > 1 && os.Args[1] == "exit" {
		shell.Process(os.Args[2:]...)
	} else {
		// start shell
		shell.Run()
		// teardown
		shell.Close()
	}
}