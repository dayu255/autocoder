# autocoder
An AtCoder CLI tool written in Golang.

> [!WARNING]
> This tool is still in development.

## How to use

- Create AtCoder workspace directory and template file.

`ac make *CONTEST_NAME* [*TEMPLATE_LANG*]`

- Show task

`ac show *CONTEST_NAME* [*TASK_LEVEL*]`

If you skip TASK_LEVEL, all level of tasks will be shown.

- Download test cases.

`ac download *CONTEST_NAME* [*TASK_LEVEL*]`

If you skil download, test cases will be downloaded at "ac test"

- Run test case

`ac test *FILE* [*CONTEST_NAME*] [*TASK_LEVEL*]`

If the file name shows the task level(Ex. a.cpp), you can skip the task level. 
If the directory name is the contest number(Ex. 440/), you can skip the contest number.

- Submit file

`ac submit *FILE* [*CONTEST_NAME*] [*TASK_LEVEL*]`

If the file name shows the task level(e.g., a.cpp), you can skip the task level. 
If the directory name is the contest number(e.g., 440/), you can skip the contest number.

