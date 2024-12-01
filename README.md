# Faraway

Это задача для технического собеседования. Она состоит из tcp-сервера и клиента. Сервер отправляет задачу клиенту. Задача заключается в вычислении proof of work. Алгоритм следующий:

Сервер генерирует случайную строку и отправляет её клиенту вместе с предустановленным префиксом. Клиент должен найти nonce на основе строки, полученной от сервера, так, чтобы hashsum строки + nonce начинался с полученного префикса. Затем клиент должен отправить nonce обратно на сервер. Сервер, в свою очередь, проверяет nonce, комбинируя его со строкой, ранее отправленной клиенту, и вычисляя hashsum комбинации. Если hashsum имеет тот же префикс, что и тот, который был отправлен клиенту, задача считается успешно выполненной.

Выбор алгоритма для вычисления proof-of-work основан на следующих пунктах:

1. Для вычисления nonce требуется время, и если префикс, отправленный сервером, длиннее, то время вычисления значительно увеличивается. Таким образом, мы контролируем сложность pow.
2. SHA-256 — это широко используемая криптографическая хеш-функция, которая генерирует хеш фиксированного размера (256 бит). Она эффективна и широко поддерживается во многих языках программирования и аппаратных средствах.
3. Выбор хеш-функции (SHA-256 в этом случае) гарантирует, что nonce не может быть легко предсказан. Это означает, что злоумышленник не может заранее вычислить nonce за разумное время без выполнения Proof of Work.
4. Этот подход прост в реализации и понимании. Он использует базовые манипуляции со строками (task + nonce), а вывод хеш-функции легко сравнивается с желаемым префиксом.