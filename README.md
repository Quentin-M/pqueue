# pqueue

[![Build Status](https://drone.io/github.com/Quentin-M/pqueue/status.png)](https://drone.io/github.com/Quentin-M/pqueue/latest)
[![Coverage Status](https://coveralls.io/repos/github/Quentin-M/pqueue/badge.svg?branch=master)](https://coveralls.io/github/Quentin-M/pqueue?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Quentin-M/pqueue)](https://goreportcard.com/report/github.com/Quentin-M/pqueue)
[![GoDoc](https://godoc.org/github.com/Quentin-M/pqueue?status.svg)](https://godoc.org/github.com/Quentin-M/pqueue)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://www.apache.org/licenses/LICENSE-2.0)

*pqueue* is an open-source collection of [priority queues] written in [Go].

[priority queues]: https://www.wikiwand.com/en/Priority_queue
[Go]: https://golang.org/

### Available structures

| Data Structure   | Push | Peek | Pop       | DecreaseKey | Has/Get | Delete    | Length | Clear |
| :------------:   | :--: | :--: | :-------: | :---------: | :-----: | :-------: | :----: | :---: |
| [Fibonacci Heap] | O(1) | O(1) | O(log n)¹ | O(1)¹       | O(1)    | O(log n)¹ | O(1)   | O(1)  |
¹ Amortized time.

[Fibonacci Heap]: https://www.wikiwand.com/en/Fibonacci_heap
