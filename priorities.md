-	Parse "" and '', so that we can have `|` inside arguments without breaking ii, for example grep 'cpu|erlang'

-	Add command blacklist : rm, mv, su, sudo, vim, vi, top, htop, nano, emacs

-	Cache output of commands when editing line at the end

-	Scrollable view for each pipe

-	Add tests (have a look at fzf tests for inspiration)

