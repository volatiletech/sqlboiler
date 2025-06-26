# Contributing

Thanks for your interest in contributing to SQLBoiler!

We have a very lightweight process and aim to keep it that way.
Read the sections for the piece you're interested in and go from
there.

If you need quick communication we're usually on [Slack](https://sqlboiler.from-the.cloud).

# New Code / Features

## Small Change

#### TLDR

1. Open PR against **master** branch with explanation
1. Participate in Github Code Review

#### Long version

For code that requires little to no discussion, please just open a pull request with some
explanation against the **master** branch. 

## Bigger Change

#### TLDR

1. Start proposal of idea in Github issue
1. After design concensus, open PR with the work against the **master** branch
1. Participate in Github Code Review

#### Long version

If however you're working on something bigger, it's usually better to check with us on the idea
before starting on a pull request, just so there's no time wasted in redoing/refactoring or being
outright rejected because the PR is at odds with the design. The best way to accomplish this is to
open an issue to discuss it. It can always start as a Slack conversation but should eventually end
up as an issue to avoid penalizing the rest of the users for not being on Slack. Once we agree on
the way to do something, then open the PR against the **master** branch and we'll commence code review
with the Github code review tools. Then it will be merged into master, and later go out in a release.

## Developer getting started

1. Add a [Configuration files](https://github.com/aarondl/sqlboiler#configuration).
1. Write your changes
1. Generate executable. Run again if you have changed anything in core code or driver code.
   ```
   ./boil.sh build all
   ```

1. Also Move sqlboiler-[driver] built to the bin of gopath if you have changed the driver code.

1. Generate your models from existing tables

   ```
   ./boil.sh gen [driver]
   ```

1. You may need to install following package before able to run the tests.

   ```
   go get -u github.com/aarondl/null
   ```

1. Test the output

   ```
   ./boil.sh test
   ```


# Bugs

Issues should be filed on Github, simply use the template provided and fill in detail. If there's
more information you feel you should give use your best judgement and add it in, the more the better.
See the section below for information on providing database schemas.

Bugs that have responses from contributors but no action from those who opened them after a time
will be closed with the comment: "Stale"

## Schemas

A database schema can help us fix generation issues very quickly. However not everyone is willing to part
with their database schema for various reasons and that's fine. Instead of providing the schema please
then provide a subset of your database (you can munge the names so as to be unrecognizable) that can
help us reproduce the problem.

_Note:_ Your schema information is included in the output from `--debug`, so be careful giving this
information out publicly on a Github issue if you're sensitive about this.
