# go-fuzz

This repo contains a simple example for [fuzzy testing in Go 1.18](https://go.dev/security/fuzz/).

## Fuzzing in Go

Fuzzing is a process where a lot of random inputs are provided to a process with the purpose to find an input to make it fail.

To use the fuzzing capabilities in the standard library write a test function that starts with `Fuzz` and accepts `f *testing.F` as the only argument.

You can then use `f.Add(values ...any)` to add inputs that will be used to generate the fuzzy values.
These could be inputs that you know are "interesting" for the function, or just inputs that you would like to test with.

Then call `f.Fuzz(func(t *testing.T, args ...any) { /* actual test goes here /* })` to write the actual test.
The values provided by the fuzzer will be in the `args` parameter.
Their types must match the ones you passed to `f.Add`.

## Examples in this repo

This repo contains a [go file](./fuzz.go) with a simple function `DontPanic(string)` that panics if called with the exact input - `"fuzz"`.
The function contains multiple nested `if` checks.
This is intentional.
When the fuzzer generates a values it determines whether that value is "interesting" or not based on change in code-coverage.

It also container the [fuzz test file](./fuzz_test.go) with the actual test - `FuzzDontPanic`.
The test just calls the function without asserting anything.
It will fail only if the function `panic`s.

## Run the example

Fuzzing is build into the Go standard library and toolchain.
To execute the fuzz test run:

```console
go test -fuzz=.
```

You should get an output like this:

```text
fuzz: elapsed: 0s, gathering baseline coverage: 0/4 completed
fuzz: elapsed: 0s, gathering baseline coverage: 4/4 completed, now fuzzing with 16 workers
fuzz: minimizing 31-byte failing input file
fuzz: elapsed: 0s, minimizing
--- FAIL: FuzzDontPanic (0.25s)
    --- FAIL: FuzzDontPanic (0.00s)
        testing.go:1356: panic: error: wrong input
            goroutine 1019 [running]:
            runtime/debug.Stack()
                /usr/local/Cellar/go/1.19.3/libexec/src/runtime/debug/stack.go:24 +0xdb
            testing.tRunner.func1()
                /usr/local/Cellar/go/1.19.3/libexec/src/testing/testing.go:1356 +0x1f2
            panic({0x11cd780, 0x123abd8})
                /usr/local/Cellar/go/1.19.3/libexec/src/runtime/panic.go:884 +0x212
            github.com/asankov/go-fuzz.DontPanic(...)
                /Users/asankov/git/asankov/go-fuzz/fuzz.go:13
            github.com/asankov/go-fuzz.FuzzDontPanic.func1(0x0?, {0xc00674bba0, 0x4})
                /Users/asankov/git/asankov/go-fuzz/fuzz_test.go:7 +0x356
            reflect.Value.call({0x11cf4e0?, 0x120b748?, 0x13?}, {0x11fd256, 0x4}, {0xc00679f500, 0x2, 0x2?})
                /usr/local/Cellar/go/1.19.3/libexec/src/reflect/value.go:584 +0x8c5
            reflect.Value.Call({0x11cf4e0?, 0x120b748?, 0x51b?}, {0xc00679f500?, 0x12ff8b0?, 0x131f5a0?})
                /usr/local/Cellar/go/1.19.3/libexec/src/reflect/value.go:368 +0xbc
            testing.(*F).Fuzz.func1.1(0x0?)
                /usr/local/Cellar/go/1.19.3/libexec/src/testing/fuzz.go:337 +0x231
            testing.tRunner(0xc0067b81a0, 0xc0067b47e0)
                /usr/local/Cellar/go/1.19.3/libexec/src/testing/testing.go:1446 +0x10b
            created by testing.(*F).Fuzz.func1
                /usr/local/Cellar/go/1.19.3/libexec/src/testing/fuzz.go:324 +0x5b9
            
    
    Failing input written to testdata/fuzz/FuzzDontPanic/a840957b7c9ea08cd6ba3e99a7ac22aa22ebecacb4656e78784d83aeb3df1189
    To re-run:
    go test -run=FuzzDontPanic/a840957b7c9ea08cd6ba3e99a7ac22aa22ebecacb4656e78784d83aeb3df1189
FAIL
exit status 1
FAIL    github.com/asankov/go-fuzz      0.373s
```

It managed to find the input which triggered the `panic`.

It also saved that input into the `testdata/fuzz/FuzzDontPanic/a840957b7c9ea08cd6ba3e99a7ac22aa22ebecacb4656e78784d83aeb3df1189` file.

That means that every time you run this test it will automatically use these values, because it already knows that they trigger a failure.
This will protect us from introducing this bug into our codebase again in the future.

To make sure this works properly checkout the [`v2`](https://github.com/asankov/go-fuzzing-101/tree/v2) branch where the bug is fixed and rerun the fuzz test to see its passing.