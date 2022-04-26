<br />
<p align="center">

  <h3 align="center">Microservice Golang (GRPC)</h3>

  <p align="center">
    Backend Blog Like Medium
    <br />
    <br />
    <a href="https://github.com/ItsMyEyes/microservices-golang/issues">Report Bug</a>
    ·
    <a href="https://github.com/ItsMyEyes/microservices-golang/issues">Request Feature</a>
  </p>
</p>

## 📎 Requirements
* [GRPC](https://grpc.io/docs/languages/go/quickstart/)
* [GO](https://go.dev/) 1.18
* [MongoDB](https://www.mongodb.com) For Database
* [Redis](https://redis.io/) For Cache


## 🚀 Installation
- Clone Repository
```
git clone https://github.com/ItsMyEyes/microservices-golang.git
```
- Install package
After cloning, go to directory and run this command for installation package
```
go mod tidy
```
- Set Enviroment
After done, set some variable or copy .env.example
```
cp .env.example .env
```
- Running development service
```
go run ./main.go
```

- Build for Production (Follow you os)
```
go build 
```

## ⚙️ Configuration (Production)
```
GO_ENV=
MONGO_URL=mongodb://<host>:<port>
REDIS_HOST=<host>:<port>
REDIS_PASSWORD=
```


## 🔐 License
Distributed under the MIT License. See [`LICENSE`](https://github.com/ItsMyEyes/microservices-golang/blob/master/LICENSE) for more information.