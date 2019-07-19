<h1 align="center">Welcome to wework-bot 👋</h1>
<p>
  <img alt="Version" src="https://img.shields.io/badge/version-v1.0.0-blue.svg?cacheSeconds=2592000" />
</p>

> tencent wework rebot

## Install

```sh
go get github.com/adnanh/webhook
for obj in `ls lib/`; do echo "build ./lib/$obj/*.go into /opt/build/$obj ..."; go build -o /opt/build/$obj ./lib/$obj/*.go; done
webhook -hooks `pwd`/hooks/hooks.yaml -verbose

```

## # Use docker
```sh
docker pull notices/wework-bot:latest
docker up -d --name wework-bot -p 9000:9000 notices/wework-bot:latest

```

## Author

👤 **shaddock**

* Github: [@ntfs32](https://github.com/ntfs32)

## Show your support

Give a ⭐️ if this project helped you!

***
_This README was generated with ❤️ by [readme-md-generator](https://github.com/kefranabg/readme-md-generator)_