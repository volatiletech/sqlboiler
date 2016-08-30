vala [![GoDoc](https://godoc.org/github.com/kat-co/vala?status.svg)](https://godoc.org/github.com/kat-co/vala)
====

A simple, extensible, library to make argument validation in Go palatable.

Instead of this:

```go
func BoringValidation(a, b, c, d, e, f, g MyType) {
  if (a == nil)
    panic("a is nil")
  if (b == nil)
    panic("b is nil")
  if (c == nil)
    panic("c is nil")
  if (d == nil)
    panic("d is nil")
  if (e == nil)
    panic("e is nil")
  if (f == nil)
    panic("f is nil")
  if (g == nil)
    panic("g is nil")
}
```

Do this:

```go
func ClearValidation(a, b, c, d, e, f, g MyType) {
  BeginValidation().Validate(
    IsNotNil(a, "a"),
    IsNotNil(b, "b"),
    IsNotNil(c, "c"),
    IsNotNil(d, "d"),
    IsNotNil(e, "e"),
    IsNotNil(f, "f"),
    IsNotNil(g, "g"),
  ).CheckAndPanic() // All values will get checked before an error is thrown!
}
```

Instead of this:

```go
func BoringValidation(a, b, c, d, e, f, g MyType) error {
  if (a == nil)
    return fmt.Errorf("a is nil")
  if (b == nil)
    return fmt.Errorf("b is nil")
  if (c == nil)
    return fmt.Errorf("c is nil")
  if (d == nil)
    return fmt.Errorf("d is nil")
  if (e == nil)
    return fmt.Errorf("e is nil")
  if (f == nil)
    return fmt.Errorf("f is nil")
  if (g == nil)
    return fmt.Errorf("g is nil")
}
```

Do this:

```go
func ClearValidation(a, b, c, d, e, f, g MyType) (err error) {
  defer func() { recover() }
  BeginValidation().Validate(
    IsNotNil(a, "a"),
    IsNotNil(b, "b"),
    IsNotNil(c, "c"),
    IsNotNil(d, "d"),
    IsNotNil(e, "e"),
    IsNotNil(f, "f"),
    IsNotNil(g, "g"),
  ).CheckSetErrorAndPanic(&err) // Return error will get set, and the function will return.

  // ...

  VeryExpensiveFunction(c, d)
}
```

Tier your validation:

```go
func ClearValidation(a, b, c MyType) (err error) {
  err = BeginValidation().Validate(
    IsNotNil(a, "a"),
    IsNotNil(b, "b"),
    IsNotNil(c, "c"),
  ).CheckAndPanic().Validate( // Panic will occur here if a, b, or c are nil.
    HasLen(a.Items, 50, "a.Items"),
    GreaterThan(b.UserCount, 0, "b.UserCount"),
    Equals(c.Name, "Vala", "c.name"),
    Not(Equals(c.FriendlyName, "Foo", "c.FriendlyName")),
  ).Check()

  if err != nil {
  	return err
  }

  // ...

  VeryExpensiveFunction(c, d)
}
```

Extend with your own validators for readability. Note that an error should always be returned so that the Not function can return a message if it passes. Unlike idiomatic Go, use the boolean to check for success.

```go
func ReportFitsRepository(report *Report, repository *Repository) Checker {
	return func() (passes bool, err error) {

		err = fmt.Errof("A %s report does not belong in a %s repository.", report.Type, repository.Type)
		passes = (repository.Type == report.Type)
		return passes, err
	}
}

func AuthorCanUpload(authorName string, repository *Repository) Checker {
	return func() (passes bool, err error) {
        err = fmt.Errof("%s does not have access to this repository.", authorName)
		passes = !repository.AuthorCanUpload(authorName)
		return passes, err
	}
}

func AuthorIsCollaborator(authorName string, report *Report) Checker {
	return func() (passes bool, err error) {

        err = fmt.Errorf("The given author was not one of the collaborators for this report.")
		for _, collaboratorName := range report.Collaborators() {
			if collaboratorName == authorName {
				passes = true
				break
			}
		}

        return passes, err
	}
}

func HandleReport(authorName string, report *Report, repository *Repository) {

	BeginValidation().Validate(
    	AuthorIsCollaborator(authorName, report),
		AuthorCanUpload(authorName, repository),
		ReportFitsRepository(report, repository),
	).CheckAndPanic()
}
```
