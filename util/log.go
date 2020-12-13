package util

import (
	"log"
	"os"
)

var Logger = log.New(os.Stdout, "[stock-checker]: ", log.Ldate|log.Ltime|log.LUTC|log.Llongfile)
