We want to create a cli that automatically removes node_modules directory by walking down specific directories and it's subdirectories.

## Features

 - Golang based
 - CLI tool
 - Only looks for node_modules directory and removes it
 - Should run on any OS
 - Should be cron friendly.
 - age based filtering, remove node_modules directory if it's older than 3 months.
 - it should accept command line arguments for the directories to walk down and the age based filtering.
