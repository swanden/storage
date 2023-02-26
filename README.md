# Storage

Написать простой GRPC сервис с командами: **get**, **set** и **delete**.
Хранилище должно быть спрятано за интерфейсом.

Реализуйте интерфейс обоими путями:

- Memcached сервер с самописной библиотекой и тремя этими же командами. [Memcached protocol](https://github.com/memcached/memcached/blob/master/doc/protocol.txt)
- Хранилище внутри памяти приложения

Оформите в виде git-репозитория и покройте тестами.

Продвинутый уровень: реализуйте пулл коннектов к memcached.

Usage
================

Run application
~~~
make up
~~~

Stop application
~~~
make down
~~~

Run integration tests
~~~
make integration-tests
~~~

Run unit tests
~~~
make unit-tests
~~~
