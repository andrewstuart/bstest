# bstest

Sometimes you just need test coverage.

`go install astuart.co/bstest@latest && bstest && go test -cover`

## Example:
```bash
➜  ~C/foo git:(master) ✗ go test -cover
?   	astuart.co/foo	[no test files]
➜  ~C/foo git:(master) ✗ bstest 99
➜  ~C/foo git:(master) ✗ go test -cover
PASS
coverage: 99.9% of statements
ok  	astuart.co/foo	0.036s
```
