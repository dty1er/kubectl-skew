// Copyright © 2021 Hidetatsu Yaginuma. All rights reserved.
package ver

import "fmt"

const escape = "\x1b"

func in(s string, c int) string {
	return fmt.Sprintf("%s[%dm%s%s[0m", escape, c, s, escape)
}

func yellow(s string) string {
	return in(s, 33)
}

func green(s string) string {
	return in(s, 32)
}

func red(s string) string {
	return in(s, 31)
}
