# ya-boo

## Пререквизиты

- OS Linux или MacOS

- установлен движок контейнеризации (docker, podman, containerd etc)

- расширение compose

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

Сохраняем api-key для дальнейшего использования
