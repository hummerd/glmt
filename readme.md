# GitLab Merge Tool

## Overview
GitLab Merge Tool (glmt) is CLI tool for making merge requests in GitLab. It's designed to be easy to use but
still flexible to cover many use cases for different teams.

## Features
* Creating MR form command line
* MR title and description with support of templates
* Team mentioning
* Slack and telegram notificationss
* View list of MR's waiting for your approval (coming soon)

## Usage

Want to crate Merge Request from command line? Just run `glmt create` from your project's directory. GLMT will
create MR from current git brunch to specified branch.

## Installation

If you have GO installed run:
```
go get gitlab.com/gitlab-merge-tool/glmt/cmd/glmt
```

## Help

Common
```
Usage:
  glmt [command]

Available Commands:
  create      Create merge request
  help        Help about any command

Flags:
  -c, --config string   path to config
  -y, --dryrun          dry run true only shows request to gitlab, but do not sends them
  -h, --help            help for glmt
  -a, --host string     gitlab host
  -l, --log string      log level (default "info")
  -k, --token string    gitlab API token

Use "glmt [command] --help" for more information about a command.
```

Create command:
```
Usage:
  glmt create [flags]

Flags:
  -d, --description string            Merge Request's description (template variables can be used in description)
  -h, --help                          help for create
  -n, --notification_message string   Additional notification message
  -b, --target string                 Merge Request's target branch (default "master")
  -t, --title string                  Merge Request's title (template variables can be used in title)
```

## Config

If you don't want to specify flags every time you can specify it in config file. By default glmt searches for glmt.config in:
* On Unix systems, it returns $XDG_CONFIG_HOME if non-empty, else $HOME/.config.
* On Darwin, in $HOME/Library/Application Support.
* On Windows, in %AppData%.
* On Plan 9, in $home/lib.

Or you can specify path to config file with `-c` flag.

Config example:
```jsonc
{
  "gitlab": { // GitLab parameters
    "url": "https://yourgitlab.com",
    "token": "XXX" // You can get one on https://YOURGITLAB.com/profile/personal_access_tokens page
  },
  "mr": { // Merge Request parameters
    "branch_regexp": "(?P<TaskType>.*)/(?P<Task>.*)/(?P<BranchDescription>.*)",
    "title": "{{.Task}} {{humanizeText .BranchDescription}}", // MR's title, can be template
    "description": "Merge feature {{.Task}} \"{{humanizeText .BranchDescription}}\" into {{.TargetBranchName}}\n{{.GitlabMentions}}", // MR's description, can be template
    "target_branch": "develop",
    "squash": true,
    "remove_source_branch": true
  },
  "notifier": { // Notification parameters
    "slack_web_hook": {
      "enabled": true,
  	  "url": "https://hooks.slack.com/services/XXX/XXX/XXX", // Learn how to get in https://api.slack.com/legacy/custom-integrations/messaging/webhooks
      "message": "<!here>\n{{.Description}}\n{{.MergeRequestURL}}" // Message to be posted to slack, can be template
    },
    "telegram": {
      "enabled": true,
      "url": "https://api.telegram.org",
      "api_key": "{{.BOT_TOKEN}}", // Ask @BotFather: https://telegram.me/BotFather.
      "message": "{{.Description}}\n{{.MergeRequestURL}}\n{{.NotificationMentions}}", // Message template.
      // Where to send a message.
      //
      // You can get group_id from:
      // https://api.telegram.org/bot{{.BOT_TOKEN}}/getUpdates.
      //
      // Also disable group privacy for the bot: https://core.telegram.org/bots#privacy-mode.
      "chat_id": "@BotFather"
    }
  },
  "mentioner": {
    "team_file_source": "PATH_TO/glmt-team.config", // Path (can be http url) to team file, see info about "Team file"
    "count": 2 // Number of project members to be mentioned in MR
  },
  // There are several environment variables available in hooks:
  // GLMT_REMOTE, GLMT_BRANCH, GLMT_MR_URL, GLMT_PROJECT, GLMT_USERNAME.
  "hooks": {
    // Before contains commands that will be executed before MR creation.
    "before": {
      "git up to date": [
        "sh", "-c", "git status | grep 'Your branch is up to date with'"
      ]
    },
    // After contains commands that will be executed after MR creation.
    "after": {
      "print mr url": [
        "sh", "-c", "echo MR: $GLMT_MR_URL"
      ]
    },
    // Timeout of all commands in the set.
    "timeout": "4s"
  }
}
```

## Templating

Title and Description and other fields can be static string or it can be template. Templates made
as https://golang.org/pkg/text/template/. In template you can specify predefined variables:
* ProjectName - project name (path extracted from git remote)
* BranchName - current branch name
* TargetBranchName - target branch name (from config or flag)
* GitlabMentions - mentions added to MR (uses username from team file, it should be gitlab username), see [Mentions](#Mentions)

Variables available for notification (previous variables are also available):
* Title - merge request title
* Description - merge request description
* MergeRequestURL - merge request URL
* NotificationMentions - mentions for notification (uses user names apropriate for this notificator), see [Mentions](#Mentions)
* ChangesCount - string with changes count for this MR

Additionally you can use any regexp group name from `branch_regexp` in description of title templates.
If `title` not specified, current branch name will be used as title. If `description` not specified
template "`Merge {{.BranchName }}" into {{.TargetBranchName}}`" will be used for description.

Also there is predefined functions for templates:
* humanizeText - capitalize first character and replaces "-" and "_" with space
* upper - change letters to upper case
* lower - change letters to lower case

## Notifications

GLMT can notify your team about created MR. Currently slack (through webhook messages) and telegram notifications are supported. Notification message also support message templating.

## Mentions

GLMT can mention some members of your team in MR and notification message. GLMT will select number of
members specified in config file (`mentioner.count`). If one of member is project owner it will be included in mention list.
All other mentioned members will be selected randomly. Only active members are mentioned.

GLMT knows your team members from team file, specified in `mentioner.team_file_source`.

### Team file

Team file has following structure:

```jsonc
{
  "members": [
    {
      "username": "john",                   // Gitlab's username (without @)
      "owns_projects": ["group/project1"],  // Project's name, owned by John
      "is_active": true,                    // Is John active at current moment (you can set it to false for vacation time)
      "names": {                            // Names for different notification channels
        "slack_member_id": "AABBXX"
      }
    },
    {
      "username": "nick",
      "owns_projects": ["group/project2"],
      "is_active": true,
      "names": {
        "slack_member_id": "CCDDXX"
      }
    }
    ...
  ]
}
```