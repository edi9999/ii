-	Parse "" and '', so that we can have `|` inside arguments without breaking ii, for example grep 'cpu|erlang'

-	Cache output of commands when editing line at the end

-	Scrollable view for each pipe

-	Add tests (have a look at fzf tests for inspiration)

Done
====

-	[20/04/2018] Add command blacklist : rm, mv, su, sudo, vim, vi, top, htop, nano, emacs
