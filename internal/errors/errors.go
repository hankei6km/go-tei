// +build go1.13

package errors

// errors.go と errors_compat.go は以下のソースコードの記事を参考に作成.
// [サポートページ：WEB&#43;DB PRESS Vol.112：｜gihyo.jp … 技術評論社](https://gihyo.jp/magazine/wdpress/archive/2019/vol112/support)
//  -「Goに入りては…… ── When In Go...」で使用されたソースコード

import "fmt"

// Wrapf wraps "%w" of fmt.Errorf on go1.13.
func Wrapf(err error, format string, a ...interface{}) error {
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), err)
}
