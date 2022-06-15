# Improve the Developer Experience with dx
`dx` is a CLI used by DevOps engineers to help generate and manage developer tooling for your company.

## The Problem
Developers at a company often need a way to work with various projects, pieces of infrastructure, containers, services, etc. and use separate tools for each.

Each tool must be installed, configured, and managed separately and even just getting everything installed can be a chore. Often there's some kind of long installation article a developer uses to get started. But the document is out of date and takes hours to complete.

Wouldn't it be nice if on Day 0 a developer could be handed a short list of commands to run on their workstation to have their environment configured?

Wouldn't it also be nice if that short list of commands introduced the developer to the company's developer experience and gave them an interactive onboarding experience as well as a long term foundation for getting their work done?

This is the problem `dx` is intended to solve.

## The Solution

DevOps platform engineers use `dx` to automatically generate and easily maintain a CLI solution for their company dev teams.

`dx` is two things.
1. A command line interface, used to generate company specific CLIs
2. A library the generated CLIs can import to help with common tasks

Generated CLIs are built based on configurable templates. The idea is that you should be able to regenerate your CLI at any time to pick up new features and bug fixes. There is a set of defaults, but you can bring your own templates as well.

## Examples
### Using dx to generate a company CLI
In this example a DevOps platform engineer at the fictious company "My Company" uses `dx` to generate a CLI for My Company's developers to use.
```
$ dx cli create myco "My Company" ~/myco
> A CLI for My Company named myco has been generated at ~/myco.
> Run $ ~/myco/bin/myco --help for more.
```

The contents of ~/myco now contains a new project you can edit, configure, and distribute to developers. You are in charge of how configuration and distribution takes place.

### Using the company CLI
Once the CLI is distributed to developer workstations, a developer can use it to perform a number of common actions such as setting up their workstation to work on a project.
```
$ myco --help
> myco is a tool for developers.
>
> Usage:
>
>    myco <command> [arguments]
>
> The commands are:
>
>    init           initialize your development environment
>    create         create a resource such as a project or virtual machine
>    list           list available resources
>    describe       describe a resource
>    destroy        destroy a resource
>    setup          configure this environment for working with an existing project
>
> Use "myco help <command>" for more information about a command.
>
> Additional help topics:
>
>    cloud          cloud providers
>    environment    environment variables
>    global         global resources available
>
> Use "myco help <topic>" for more information about that topic.
>
$ myco init
> initializing the environment...
> installing dependencies...
> tfenv is not installed. Install it? [Y/n]
> installing tfenv...
> installing terraform...
> configuring cloud providers...
> done.
>
$ myco setup ui
> cloning repositories...
> wrangling config files...
> performing custom actions...
> done.
>
> The project ui has been set up and is ready for you to work with at ~/myco/ui
$
```

