# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`modernc.org/memory` is a single-package, dependency-light Go library implementing a C-style `malloc`/`free` allocator on top of raw OS mappings (`mmap` on unix, `VirtualAlloc` on Windows). It exists to serve the modernc C-to-Go stack (`modernc.org/libc` and everything above it), which needs memory the Go GC neither moves nor scans, and addresses that stay valid when held only as `uintptr`.

The repo is flat: ten Go files, one test file. `internal/autogen/` is an empty placeholder for the modernc builder infrastructure described by `builder.json`.

## Commands

```sh
go test                                 # full suite, ~2m15s (each test churns a 128 MiB quota)
go test -tags=memory.counters           # same, with the accounting assertions actually live — see Build tags
go test -run TestUMalloc                # a single test
go test -run @ -bench . -benchmem       # benchmarks only ("@" matches no test name)
make editor                             # gofmt -l -s -w *.go — the pre-commit step
make build_all_targets                  # go build + staticcheck across every supported GOOS/GOARCH
make cpu | make mem                     # bench under a profiler, then open pprof
make todo                               # grep TODO/BUG/println/unused-assignment markers
```

`make all` is the author's personal pipeline: it needs `misspell`, `unconvert`, `maligned` and `golint` installed and greps a file named `log` that it does not create (the convention is `go test ... |& tee log`). It usually fails locally — reach for the individual targets instead.

**`go vet`**: the standalone command reports dozens of `possible misuse of unsafe.Pointer` across `memory.go` and `all_test.go` and exits 1. These are inherent to the design — the allocator deliberately round-trips addresses through `uintptr` into memory the GC does not manage — and must not be "fixed". For a meaningful signal use `go vet -unsafeptr=false ./...` or `staticcheck`; both are clean.

## Build tags

- **`memory.counters`** (`counters.go` / `nocounters.go`) — enables the `Allocator.Allocs/Bytes/Mmaps` bookkeeping. Without it those fields stay 0, which silently makes the `alloc.Allocs != 0 || alloc.Mmaps != 0 || alloc.Bytes != 0` assertion ending nearly every test vacuous; only the `len(alloc.regs) != 0` leak check survives. **Any change to mmap/unmap accounting must be tested with `-tags=memory.counters`.**
- **`memory.trace`** (`trace_enabled.go` / `trace_disabled.go`) — logs every Malloc/Free/Realloc/UsableSize call to stderr.

Both constants are compile-time `false` by default, so the guarded code costs nothing in normal builds.

## Architecture

Everything hangs off one invariant: **every mapping is `pageSize`-aligned (64 KiB, `pageSizeLog = 16`), so any allocated address masks down to its page header with `p &^ pageMask`.** There are no per-allocation headers and no address→metadata lookup on the hot path — `Free` and `UsableSize` recover everything from that single AND.

`page` (`memory.go:191`) is the four-word header at the start of every mapping: `brk` (bump index), `class` (size class index), `size` (bytes actually mapped), `used` (live slots). `init()` asserts `sizeof(page) % mallocAllign == 0`.

Two allocation paths, discriminated by `page.class`:

- **`class != 0` — shared page**, for `size <= maxSlotSize` (`1 << (pageSizeLog-2)` = 16 KiB). The request is served from the smallest size class that fits: `class = sizeToClass[(size-1)/mallocAllign]`, slot size `classes[class].size`. A page holds `classes[class].slots = pageAvail / size` slots handed out by bumping `brk`; `a.pages[class]` is the current carve target and is zeroed once it fills. Freed slots go on `a.lists[class]`, an intrusive doubly-linked list whose `node{prev,next}` is written *inside the free slot* — hence the two-word (`mallocAllign`) minimum slot size.
- **`class == 0` — dedicated page**, for anything larger. One allocation per mapping at `pg + headerSize`, released on `Free` (see Page retention). Class indexes start at 1, so 0 is free to be the sentinel.

**Size classes** (`initClasses`, `memory.go:147`): `linearClasses` (8) multiples of `mallocAllign`, then `classesPerDoubling` (4) evenly spaced classes per power of two, up to `maxSlotSize` — on 64-bit 16, 32, …, 128, 160, 192, 224, 256, 320, …, 16384 (36 classes; 40 on 32-bit, where the quantum is 8). Every power of two stays a class. Slots sit back to back after the 32-byte header, so a slot is `mallocAllign`-aligned and no more — finer alignment must not be assumed. The tables are generated at init; `initClasses(1, 1)` reproduces the pre-2026-11 power-of-two classes exactly, which is how `churn_test.go` A/B-compares the two layouts in one binary. Class indexes are masked with `maxClasses-1` (a power of two) to drop bounds checks on the hot path. modernc.org/libc's tcache (`mem_musl.go`) still rounds requests ≤ 1 KiB to powers of two and caches only power-of-two usable sizes; it relies on every power of two being a class. `UsableSizeFor` exists so that it can bin by the real classes instead of copying the layout, which is how its `tcacheLog` comment ("as modernc.org/memory picks it") went stale. The class-binned tcache (prototyped in libc 2026-11-15, to land after the memory tag) builds its bins from `UsableSizeFor` at init and panics there unless: `UsableSizeFor(n)` is in `[n, 1024]` for every `n <= 1024` (1024 stays a class); there are at most 32 classes up to 1 KiB; and `UsableSizeFor(s) == s` for every class size `s`. Class tuning must keep all three.

`a.regs` is the set of live mappings. It is touched only by `mmap`/`unmap`/`Close` (and by tests as a leak check), never on the allocation fast path.

**Page draining** (`UintptrFree`, `memory.go:476-493`): when a shared page's `used` reaches 0, all `brk` of its slots are unlinked from the global `a.lists[class]` before the page is retained or unmapped. Skipping that walk would leave free-list nodes pointing into memory that is no longer on the list — or no longer mapped at all.

**Alignment costs syscalls** (`mmap_unix.go:34`): to get a 64 KiB-aligned result, `mmap` asks for `size + pageSize` and unmaps the misaligned head and the surplus tail — up to three syscalls per page acquisition. `Allocator.mmap` amortizes this for 64 KiB regions by acquiring `slabBatch` of them in one mapping and pooling the rest (unix only — `canCarve`). If the *tail* unmap is rejected by the kernel the code keeps the enlarged mapping rather than failing (fixes bigsort on linux/s390x, cznic/sqlite#207), so `page.size` may legitimately exceed what was requested. The `TODO` at `memory.go:314` is about reusing that surplus, which requires moving the slot count out of the global class table and into `page` so capacity becomes per-page.

Windows gets alignment for free — `VirtualAlloc` is 64 KiB-granular — so `mmap_windows.go` just rounds up and returns.

**Platform and width splits**: `pageSizeLog` is declared per platform in `mmap_unix.go` and `mmap_windows.go` (16 in both since `e95c668`). `rawmem`, the huge array type used to fabricate slices over raw memory, is declared per pointer width in `memory32.go` / `memory64.go` — its build-tag list is what needs editing when adding an architecture.

## API surface

Four families over one core, all in `memory.go`: `Malloc`/`Calloc`/`Realloc`/`Free` returning `[]byte`, the `Unsafe*` variants returning `unsafe.Pointer`, the `Uintptr*` variants, and `UsableSize`/`UnsafeUsableSize`/`UintptrUsableSize`. **The `Uintptr*` methods are the implementation; the other two families are thin wrappers** — algorithm changes belong in `UintptrMalloc`/`UintptrFree`/`UintptrRealloc` and nowhere else.

`UsableSizeFor(size)` (package-level, after `UsableSize`) is the `nallocx(3)` of this package: the usable size `Malloc(size)` produces, computed from the class tables alone — exact up to `maxSlotSize`, `mmapSize(size+headerSize)-headerSize` above (a block exceeds that only when the kernel refused to trim the tail of a fresh mapping, see Alignment costs syscalls). It reads no `Allocator`, so it is safe for concurrent use. It does not inline — the out-of-line call for big sizes alone costs 57 of the 80-node inline budget — and costs about 2.4 ns against 0.5 ns for a table lookup, so a hot caller should build its own table from it at init. Any change to how `UintptrMalloc` or `newPage` size a block must keep it in step; `TestUsableSizeFor` checks it against real blocks of every slot size in both class layouts.

The zero `Allocator` is ready to use and **is not goroutine-safe** — there is no locking anywhere; callers serialize. Slice-returning calls set cap to `usableSize(p)`, which is why `Free` and `Realloc` reslice to `b[:cap(b)]` to recover the block address.

## Tests

`all_test.go` is the entire suite: randomized allocate / verify / free cycles over a 128 MiB quota, in `Small` (≤ 2× OS page, shared pages only) and `Big` (≤ 2× 64 KiB, forces the dedicated-page path) variants, across all three access styles. Sizes and contents come from `mathutil.NewFC32`, a full-cycle PRNG the tests `Seek(0)` to replay — expected values are regenerated rather than stored, and failures reproduce deterministically. Every allocation is filled with its pattern and re-verified, so slot overlap or a wrong `usableSize` surfaces as "corrupted heap". Tests read unexported state (`alloc.regs`, `page.class`), so they must stay in-package.

`churn_test.go` adds a steady-state mixed-size workload: `BenchmarkChurn` (ns per free+malloc) and `TestFootprint` (usable, mapped and — on linux, via `mincore` — resident bytes against requested), each over several size distributions, including `qjsSizes`, the 64 most frequent request sizes QuickJS makes on the V8 suite, and over both class layouts (`classSchemes`).

The `Malloc`/`Free`/`Calloc` benchmarks in `all_test.go` allocate arrays of `b.N` blocks, so the default benchtime drives them to several GB; on a shared host pass a fixed `-benchtime=10000000x`.

## Page retention

`UintptrFree` does **not** return an empty region to the OS eagerly. A drained shared page is kept as its size class's carve target; every other empty region — a drained page whose class already has a carve target, or a dedicated page — goes via `release` into `a.freed`, a pool keyed by region size and bounded by `maxFreedSize` (4 MiB plus the high-water mark of the live mapping). `Allocator.mmap` serves requests from the pool before mapping fresh, with request sizes collapsed into region classes (`mmapSize`) so pool hits are likely. `Trim` releases the carve targets and the pool; `Close` releases everything.

This is why every test and benchmark calls `alloc.Trim()` before asserting `Allocs`/`Mmaps`/`Bytes`/`regs` are zero. Asserting after `Close` instead would be vacuous — `Close` ends with `*a = Allocator{}`, so every field reads zero regardless.

`HANDOFF.md` records the diagnosis, the measurements, and the design options that were weighed before picking this policy. Its "Status: open" header is historical; note also that it undercounts the affected test sites as 8 (there are 16 — 10 using `t`, 6 using `b`).
