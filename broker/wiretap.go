package broker

import (
	"PAD-151-Message-Broker/model"
	"log"
	"os"
)

// WireTap - referes to working with files
type WireTap struct {
	file *os.File
	err  error
}

// Init - initiates file
func (wt *WireTap) Init(fileName string) {
	wt.file, wt.err = os.Create("test.txt")
	if wt.err != nil {
		log.Fatal(wt.err)
	}
}

func (wt *WireTap) Append(sentMessageModel *model.SentMessageModel) {
	data, err := model.EncodeJsonMessage(*sentMessageModel)
	if err != nil {
		log.Println(err)
	}
	wt.file.Write([]byte(append(data, "\n"...)))
}

// Close - closes create file
func (wt *WireTap) Close() {
	wt.file.Close()
}
