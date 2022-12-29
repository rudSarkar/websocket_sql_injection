# WebSocket Sql Injection

I came accross a HackTheBox machine called [Soccer](https://app.hackthebox.com/machines/Soccer) which is pretty interesting in a point of WebSocket SQL Injection, Did some reaseach and found a blog from [salmonsec](https://salmonsec.com/cheatsheet/web_sockets) and it's written in Python2, as we know Python2 is deprecated. I used golang and made this to easier my work. The specialty of my tool it can take multiple parameter eg: http://localhost:8000/?key1=value1&key2=value2 

#### Installation

```bash
go install github.com/rudSarkar/websocket_sqli@latest
```

#### Server

Copy the WebSocket URL and use the `-ws` flag to pass the URL of WebSocket it will up and running then follow the SQLMAP

```bash
websocket_sqli -ws 'ws://soc-player.soccer.htb:9091/'
```

#### SQLMAP

```bash
sqlmap -u "http://127.0.0.1:8000/?id=73456" -p id --random-agent --dbs
```

