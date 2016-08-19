# ps2bs

    go get github.com/codemodus/ps2bs

ps2bs (pointers to built-ins) is a CLI application for generating helper
functions which return pointers to built-ins and other commonly used standard
library types.

The code is formatted and satisfies all common linters.

    Available flags:

    --dir={dir}      Set the destination for the generated file.
    -e               Set the exportation of functions.

This command is not meant to be used with go:generate. It's purpose is to
provide enough convenience that importing a new dependency into a project is
less appealing than running this command and placing the generated code under
version control.

The following usage will produce helpers within the current directory's Go
package and store them in a file named "./ps2bs_gen.go".

    ps2bs

No concern need be given for the package name as it will be analyzed and used.
If it is desirable for the generated code to be in a sub-package and accessed
from the current directory or others, use the following:

    // where "subdir" is an existing directory
    ps2bs -e --dir=subdir

In this case, all functions will be exported and the package name will be
derived from the dir flag's value.

If a sub-directory already exists with other Go code within, the previous
command will work as before, but the package name will be aligned with the
existing code. In the event that Go naming conventions are violated or the code
in the destination directory is not buildable, the command will exit with a
non-zero status and provide an error message.
