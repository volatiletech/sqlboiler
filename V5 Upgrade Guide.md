# V5 Upgrade Guide

## CLI Changes

* The `--no-context` flag has been removed. All models are now generated using context.  
    Similarly, interfaces and methods for non-context variants have been removed.
    e.g. `boil.Executor` vs `boil.ContextExecutor`.
* The `--add-enum-type` flag has been removed. Generating types for enums is now the default behavior.
* Removed the `--no-rows-affected` flag. Methods will now always return the rows affected.

## Fixes that caused breaking changes

* Both `qm.Limit` and `qm.Offset` now use `int64` instead of `int`. You may have to type cast your current calls.
* `var TableNames` in the generated models are now based on the table alias and not a title-cased version of the table name.
* Default type for custom nullable postgres types is now `null.String` instead of `string`

