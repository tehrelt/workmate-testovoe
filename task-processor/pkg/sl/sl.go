package sl

import "log/slog"

func Err(e error) slog.Attr {
	return slog.String("err", e.Error())
}
