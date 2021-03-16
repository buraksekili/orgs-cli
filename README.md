# orgs-cli

CLI of https://organizationtechs.netlify.app/

It fetches repositories of organizations on GitHub and sorts them based on star counts.

## Installation

```
$ git clone https://github.com/buraksekili/orgs-cli.git
$ cd orgs-cli
$ go build
$ ./orgs-cli go turkey
```

## GitHub Token

You need GitHub access token to run program. 

Check https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token

You can specify your token in `.env` file.
```
GITHUB_TOKEN=yourtoken
```


![orgs-cli](https://user-images.githubusercontent.com/32663655/110846394-c675bf00-82bc-11eb-88cc-2f97bdcafd65.png)

