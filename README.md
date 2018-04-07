ii
==

*Run pipes interactively, without EVER hitting enter*

**ii** , standing for "**i**nteractive **i**nteractive", is a new way of writing multiple piped commands and get instant visual feedback.

For example, if I want to write the bash oneliner :

ps -axf | grep cpu | grep -v cpuhp | sed 's/--.*//g' | grep -o '/.*' | grep -v erlang | uniq

It is cool to have a way to see automatically what is happening.

Here's a demo :

![demo gif](https://raw.github.com/edi9999/i/master/demo_ii.gif?v=2)
