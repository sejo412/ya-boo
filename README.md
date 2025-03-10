# ya-boo

## Пререквизиты

- OS Linux или MacOS

- установлен движок контейнеризации (docker, podman, containerd etc)

- расширение compose для движка

## Подготовка инфраструктуры

**Здесь и далее подразумевается запуск LLM на локальной рабочей станции разработчика (чтобы запустилось и как-то работало).** Для запуска на GPU или тонкой настройки обращайтесь к [документации](https://github.com/ggerganov/llama.cpp)

Выбираем модель на [https://huggingface.co/](https://huggingface.co/). Для разработки/тестирования рекомендуется [https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF](https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF)

```
$ mkdir -p llm/models
$ wget -O llm/models/dev-model.gguf https://huggingface.co/bartowski/Meta-Llama-3.1-8B-Instruct-GGUF/resolve/main/Meta-Llama-3.1-8B-Instruct-Q4_K_M.gguf
```

Запуск сервера (если падает, тюним движок контейнеризации - на MacOS потребовалось увеличение выделяемой оперативной памяти):

```
$ podman run -v $(pwd)/llm/models:/models -p 8000:8000 ghcr.io/ggerganov/llama.cpp:server -m /models/dev-model.gguf --port 8000 --host 0.0.0.0 --ctx-size 100
```

Регистрируем бота в телеграм [https://telegram.me/BotFather](https://telegram.me/BotFather):

```
/start
/newbot
<Видимое имя>
<имя_бота>_bot
```

Сохраняем api-key (далее TGKEY) для дальнейшего использования

## Конфигурация и запуск

Конфигурация поддерживает конфигурационный файл, переменные окружения, ключи запуска (в порядке увеличения приоритета)

config.yaml:

```
dsn: postgres://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable  # строка подключения к БД
tgSecret: <TGKEY>                                                           # полученный ранее api-key для telegram бота
initBotSecret: xxxxxx                                                       # парольная фраза для инициализации бота (регистрация первого администратора бота)
```

Для запуска *all-in* можно воспользоваться расширением *compose* для движков контейнерезации, для этого необходимо создать минимальный *.env* файл в корне проекта:

```
$ echo BOO_TGSECRET=<TGKEY> > .env            # api-key для telegram-бота
$ echo BOO_INITBOTSECRET=<SECRET> >> .env     # парольная фраза для инициализации бота
```

Запуск:

```
$ make up || task up
```

Остановка:

```
$ make down || task down
```

Добавляем бота себе в контакт-лист `https://t.me/<имя_бота>_bot`

При первом запуске бот находится в режими инициализации, ожидая регистрации первого администратора. Для запуска и регистрации себя в качестве первого администратора необходимо знать парольную фразу, заданную при конфигурировании (*BOO_INITBOTSECRET*):

```
/init <SECRET>
```

Регистрация языковой модели на примере скачанной ранее в разделе *Подготовка инфраструктуры*:

```
/llmadd name=local endpoint=http://llama:8000 desc=локальная_модель_meta
```

Регистрация внешней LLM на примере ChatGPT (**Внимание! Подписка ChatGPT PLUS/PRO не означает доступ к API**, api-запросы тарифицируются отдельно)

```
/llmadd name=local endpoint=https://api.openai.com token=<token> desc=chatGPT
```

*TODO* [доработать парсилку с поддержкой пробелов в значениях параметров](https://github.com/sejo412/ya-boo/issues/20)

Вывод зарегистрированных моделей:

```
/llmlist
```

*TODO* [BUG теряется поле description при добавлении llm](https://github.com/sejo412/ya-boo/issues/21)

*TODO* [Пофиксить разметку Markdown](https://github.com/sejo412/ya-boo/issues/14)

Переключение на добавленную модель:

```
/llmuse 1
```

Дальнейшие запросы, не начинающиеся на "/" начнут обрабатываться выбранной моделью

При добавлении бота в контакт-лист пользователи telegram регистрируются с ролью Unknown и ждут аппрува администратора:

```
/list
/approve <ID>
```

## Команды администратора

```
/init
/approve
/list
/ban
/llmadd
/llmremove
/llmlist
/llmuse
```

## Команды пользователя

```
/llmlist
/llmuse
```

## Известные баги/туду

[Issues](https://github.com/sejo412/ya-boo/issues)
