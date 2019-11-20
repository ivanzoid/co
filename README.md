# co

Save your time by including ticket number in git commit message automatically. For those who prefers to commit from terminal.

# Usage

- Assuming your branch name is `feature/name/ID-42`,

- You may run the following command:  
`co Hello this is commit message` 
  
  which will be converted to (and executed):  
`git -a commit -m "ID-42  Hello this is commit message"`
  
- For this to work you should create `~/.co/config` with following contents:  
  `\w*\/\w*\/[\w-]*\.([-_a-zA-Z0-9]*)`  
  
  This is regexp which tries to match ticket number from branch name. Feel free to change it to match your project's branch naming conventions.
