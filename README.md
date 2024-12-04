# Create a program to send Event Viewer errors & potentially other information to help troubleshoot PC issues

(Nephew is having PC issues, he's 13 and live in New Zealand while I live in Australia, and I'd love to help him resolve them)

## How to use

- install & run postgres db

- create a .env file in the base directory
- fill out the below variables

```
DB_URL=<dburl>
SECRET=<secret>
PORT=<port>
```

- update the address inside event_logger.ps1 to point to your server

- run main.go, run event_logger.ps1 on the computer you want logs sent from
