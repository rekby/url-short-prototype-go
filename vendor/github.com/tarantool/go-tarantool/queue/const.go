package queue

const (
	READY   = "r"
	TAKEN   = "t"
	DONE    = "-"
	BURIED  = "!"
	DELAYED = "~"
)

type queueType string

const (
	FIFO      queueType = "fifo"
	FIFO_TTL  queueType = "fifottl"
	UTUBE     queueType = "utube"
	UTUBE_TTL queueType = "utubettl"
)
