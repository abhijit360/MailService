# Mail Service for my personal contact me form

I built this service that serves an API at a fixed api that I can responses to my contact me form on my website instead of using an existing tool. This tool uses gmail's smtp server to have an automated response to the client stating that I have received their email and sending my personal email a reminder to reach out to the client.

### Set up:
Clone the repo and create a set up an .env file into the project.

Your env file should look like:

```
my_mail= Username
my_password= Password (create an APP password. Not your email's password)
```

You should also change the variable for salutationName to sign off by a different name/moniker