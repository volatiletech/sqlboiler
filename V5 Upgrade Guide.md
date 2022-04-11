# V5 Upgrade Guide

## Fixes that caused breaking changes

* Both `qm.Limit` and `qm.Offset` now use `int64` instead of `int`. You may have to type cast your current calls.
* `var TableNames` in the generated models are now based on the table alias and not a title-cased version of the table name.
* Default type for custom nullable postgres types is now `null.String` instead of `string`



