/* Initial source from: https://github.com/bogem/nehm/blob/master/logs/logs.go
   Under MIT license: https://github.com/bogem/nehm/blob/master/LICENSE
*/

package logs

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var (
	INFO     = log.New(ioutil.Discard, "", 0)
	WARN     = log.New(os.Stdout, YellowString("WARN: "), 0)
	ERROR    = log.New(os.Stderr, RedString("ERROR: "), 0)
	FATAL    = log.New(os.Stderr, RedString("FATAL ERROR: "), 0)
	FEEDBACK = new(feedback)
)

func EnableInfo() {
	INFO = log.New(os.Stdout, "INFO: ", 0)
}

type feedback struct{}

func (feedback) Print(a ...interface{}) {
	fmt.Print(a...)
}

func (feedback) Println(a ...interface{}) {
	fmt.Println(a...)
}

func (feedback) Printf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}
