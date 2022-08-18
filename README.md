[![Commitizen friendly](https://img.shields.io/badge/commitizen-friendly-brightgreen.svg)](http://commitizen.github.io/cz-cli/)

# SOVA Caller App

## Разработка

### Требования

- Go 1.18.3 (https://go.dev/doc/install)
- Node v14.18 (https://nodejs.org/en/download/)
- npm latest
- PJSIP v.2.12 ([Загрузите PJSIP](https://www.pjsip.org/download.htm) и [Соберите и установите](https://trac.pjsip.org/repos/wiki/Getting-Started)). При необходимости установите отсутствующие библиотеки (libssl-dev, uuid-dev, libasound2-dev)

### Установка зависимостей

Для установки зависимостей проекта в корневой директории sova-caller выполните команды:

```
npm i
cd ./ui
npm i
cd ./server
npm i
```

Для запуска проекта в режиме разработки в корневой директории выполните:

```
npm run start
```

Бекенд-часть поднимется локально на 4000 порту [http://localhost:4000](http://localhost:4000)
UI-часть запустится и будет доступна в браузере по адресу [http://localhost:3000](http://localhost:3000)

## Сборка образа и запуск контейнера

### Требования

- Node v14.18
- npm latest
- Docker latest

В образе происходит сборка обоих частей проекта, он содержит необходимые зависимости (Go и PJSIP), поэтому для запуска приложения из контейнера достаточно выполнить только указанные ниже команды.

Для установки зависимостей проекта в корневой директории выполните команды:

```
npm i
```

Для сборки образа из Dockerfile выполните:

```
docker build -f ./docker/Dockerfile -t sova-caller .
```

Для запуска контейнера с присвоением ему имени my-sova-caller выполните:

```
docker run --name my-sova-caller -p 4000:4000 sova-caller
```

Приложение поднимается локально на 4000 порту [http://localhost:4000](http://localhost:4000)

Для остановки и запуска контейнера можно выполнить следующие команды:

```
docker stop my-sova-caller
docker start my-sova-caller
```
