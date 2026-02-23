# Go Concurrency Patterns

A structured collection of core Go concurrency patterns focused on channel-based coordination and goroutine orchestration.

This repository demonstrates practical implementations of streaming, fan-in/fan-out composition, worker pools, backpressure control, rate limiting, and channel orchestration patterns.

All patterns are implemented using Go generics.

Each pattern is written in an idiomatic, minimal, and race-safe way.

---

## Patterns

- [x] Generator
- [x] Fan-out
- [x] Fan-in
- [x] Fan-in + Fan-out
- [ ] Pipeline
- [ ] Worker Pool
- [ ] Bounded Worker Pool
- [ ] Semaphore (channel-based)
- [x] Rate Limiting
- [ ] Backpressure
- [ ] Tee Channel
- [ ] Bridge Channel