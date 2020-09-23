# WIP: GitLab Merge Tool

## Overview
GitLab Merge Tool (glmt) is CLI tool for making merge requests in GitLab. It's designed to be easy to use but still flexible to cover many use cases for different teams.

## Features
* Creating MR form command line
* MR title and description with support of templates
* Team mentioning
* View list of MR's waiting for your approval

## Usage
```
Usage:
  glmt [command]

Available Commands:
  create      Create merge request
  help        Help about any command

Flags:
  -c, --config string   path to config
  -h, --help            help for glmt
  -k, --token string    gitlab API token

Use "glmt [command] --help" for more information about a command.
```