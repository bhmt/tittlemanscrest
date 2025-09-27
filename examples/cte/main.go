package main

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/bhmt/tittlemanscrest/api"
	"github.com/bhmt/tittlemanscrest/api/handlers"
	"github.com/bhmt/tittlemanscrest/cmd"
)

var output = []string{
	"My name is Jon Daker.\n",
	"'Ord is risen today\n",
	"Aaahhhhggglayloooya\n",
	"Suts of man and angels sayy\n",
	"Aaahhhhggglayloooya\n",
	"'ts your voicean triumphs fooooiiii\n",
	"Aaahhhhggglayloooya\n",
	"sn' boyce in taoooowl\n",
	"Aaahhhhhhhhhhhooooooouulgh!\n",
	"When the\n",
	"moon hits your eye like a big preetza pie\n",
	"that's amoreee.\n",
	"Whep?\n",
	"...yrgn...\n",
	"samoreee\n",
	"Bells will ring, ting-a-la-ling\n",
	"tingawlinnnng as a bell\n",
	"bing\n",
	"mawrayyyee.\n",
	"..ng...\n",
	"ticketickey tay\n",
	"t'say hmm\n",
	"mawrayyyee.\n",
	"..aloopa...\n",
	"scrum screeee\n",
	"mmmmmmm\n",
	"Mr. Moraaayyee\n",
	"mmmmmmmmmm\n",
	"oosh shines hmmm\n",
	"dits\n",
	"You're in love\n",
	"when you know\n",
	"that you fnnn\n",
	"nnnsome more ayyeeerr\n",
	"s'me but you see back in old Napoli\n",
	"that's amoreee.\n",
}

type Loop struct {
	row, col int
}

func (l *Loop) Read(p []byte) (int, error) {
	time.Sleep(1 * time.Second)
	line := []byte(output[l.row])
	chars := line[l.col:]

	if len(chars) <= len(p) {
		l.row = (l.row + 1) % len(output)
		l.col = 0
	} else {
		l.col += len(p)
	}

	n := copy(p, chars)
	return n, nil
}

func reader() io.Reader {
	return &Loop{}
}

func Work(ctx context.Context) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)).With(slog.String("app_id", "cte"))

	mux := http.NewServeMux()
	mux.Handle(
		"/cte",
		api.MiddlewareBase(
			logger,
			http.HandlerFunc(handlers.ChunkedTransferEncoding(reader)),
		),
	)

	server := api.New(":8081", mux)
	go func() { server.ListenAndServe() }()
	logger.InfoContext(ctx, "listening on :8081")

	<-ctx.Done()
}

func main() {
	cmd.Run(Work)
}
