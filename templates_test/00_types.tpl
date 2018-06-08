var (
    // Relationships sometimes use the reflection helper queries.Equal/queries.Assign
    // so force a package dependency in case they don't.
    _ = queries.Equal
)
