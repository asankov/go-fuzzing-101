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

## Fix the code and rerun the example

We [fixed the code](https://github.com/asankov/go-fuzzing-101/commit/cf9dda299d05bd20fd29684ad1f4d38cb2067adf) to not `panic` in this case.

Let's rerun the fuzz test:

```console
$ go test -fuzz=.
fuzz: elapsed: 0s, gathering baseline coverage: 0/6 completed
fuzz: elapsed: 0s, gathering baseline coverage: 6/6 completed, now fuzzing with 16 workers
fuzz: elapsed: 3s, execs: 763627 (254535/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 6s, execs: 1575414 (270595/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 9s, execs: 2434797 (286379/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 12s, execs: 3298148 (287863/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 15s, execs: 4104638 (268803/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 18s, execs: 4895829 (263679/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 21s, execs: 5677917 (260749/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 24s, execs: 6446923 (256357/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 27s, execs: 7208629 (253908/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 30s, execs: 7965502 (252212/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 33s, execs: 8720265 (251663/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 36s, execs: 9467390 (249048/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 39s, execs: 10211900 (248089/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 42s, execs: 10958031 (248776/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 45s, execs: 11688527 (243505/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 48s, execs: 12434381 (248626/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 51s, execs: 13160188 (241937/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 54s, execs: 13874856 (238220/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 57s, execs: 14581058 (235396/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m0s, execs: 15266795 (228586/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m3s, execs: 15919550 (217585/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m6s, execs: 16607432 (229293/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m9s, execs: 17314867 (235814/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m12s, execs: 17917331 (200770/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m15s, execs: 18536553 (206451/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m18s, execs: 19252068 (238505/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m21s, execs: 19957670 (235132/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m24s, execs: 20636045 (226174/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m27s, execs: 21274409 (212811/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m30s, execs: 21981651 (235720/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m33s, execs: 22706097 (241496/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m36s, execs: 23377195 (223712/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 1m39s, execs: 24077053 (233287/sec), new interesting: 0 (total: 6)
fuzz: elapsed: 0s, gathering baseline coverage: 0/6 completed
fuzz: elapsed: 0s, gathering baseline coverage: 6/6 completed, now fuzzing with 16 workers
PASS
ok      github.com/asankov/go-fuzz      102.707s
```

We notice a few things:

1. It managed to cover the whole file pretty quick:

    ```text
    fuzz: elapsed: 0s, gathering baseline coverage: 0/6 completed
    fuzz: elapsed: 0s, gathering baseline coverage: 6/6 completed, now fuzzing with 16 workers
    ```

    That is because we already have the "interesting" input in the `testdata` folder.
    The fuzz test used that and hot full coverage of the file.

    However, this time this did not triggered a fail.

2. It continued to produce new inputs

    Since it didn't fail with the previously bad input it continued to generate new ones, until we stopped it manually:

    ```text
    fuzz: elapsed: 0s, gathering baseline coverage: 0/6 completed
    fuzz: elapsed: 0s, gathering baseline coverage: 6/6 completed, now fuzzing with 16 workers
    ```

    By default fuzz test run indefinetely until they find a wrong input.
    That is why it's a good idea to provide a timeout to the fuzz command:

    ```console
    $ go test -fuzz=. -fuzztime=30s
    fuzz: elapsed: 0s, gathering baseline coverage: 0/6 completed
    fuzz: elapsed: 0s, gathering baseline coverage: 6/6 completed, now fuzzing with 16 workers
    fuzz: elapsed: 3s, execs: 723846 (241279/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 6s, execs: 1571657 (282596/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 9s, execs: 2385118 (271094/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 12s, execs: 3153404 (256147/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 15s, execs: 3947954 (264795/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 18s, execs: 4694179 (248784/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 21s, execs: 5461521 (255791/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 24s, execs: 6219059 (252487/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 27s, execs: 7003528 (261498/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 30s, execs: 7735617 (243810/sec), new interesting: 0 (total: 6)
    fuzz: elapsed: 30s, execs: 7735617 (0/sec), new interesting: 0 (total: 6)
    PASS
    ok      github.com/asankov/go-fuzz      30.220s
    ```

    Now, the test run for 30 seconds, and since it was not able to find a failing input for that time it was marked as passing.

## References

These example are based on two great talks I watched on the subject:

- [Fuzzing in Go by Valentin Deleplace, Devoxx Belgium 2022](https://www.youtube.com/watch?v=Zlf3s4EjnFU)
- [Write applications faster and securely with Go by Cody Oss, Go Day 2022](https://www.youtube.com/watch?v=aw7lFSFGKZs)

Check them out if you want to learn more about fuzzing.
