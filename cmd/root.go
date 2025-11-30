/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreServices

#include <CoreServices/CoreServices.h>
#include <stdlib.h>

char* getDefinition(const char* word) {
	CFStringRef wordRef = CFStringCreateWithCString(NULL, word, kCFStringEncodingUTF8);
	CFRange range = CFRangeMake(0, CFStringGetLength(wordRef));

	// Get definition from system dictionaries
	CFStringRef definition = DCSCopyTextDefinition(NULL, wordRef, range);

	CFRelease(wordRef);

	if (definition == NULL) {
		return NULL;
	}

	// Convert CFString to C string
	CFIndex length = CFStringGetLength(definition);
	CFIndex maxSize = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
	char *buffer = (char *)malloc(maxSize);

	if (CFStringGetCString(definition, buffer, maxSize, kCFStringEncodingUTF8)) {
	CFRelease(definition);
		return buffer;
	}

	free(buffer);
	CFRelease(definition);
	return NULL;
}
*/
import "C"
import (
	"os"
	"runtime"
	"strings"
	"unsafe"

	"github.com/J-Siu/go-dict/global"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-dict",
	Version: global.Version,
	Short:   "MacOS command line dictionary",
	Long:    "MacOS command line dictionary. Due to API limitation, only New Oxford American Dictionary is used.",
	Run: func(cmd *cobra.Command, args []string) {
		debug, _ := cmd.Root().Flags().GetBool("debug")
		if debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		for _, word := range args {
			def := definition(word)
			if debug {
				ezlog.Debug().N("Definition").M(def).Out()
			}
			output(def)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	switch runtime.GOOS {
	case "darwin":
	default:
		ezlog.Log().M("Only support MacOS.").Out()
		os.Exit(1)
	}

	rootCmd.Flags().BoolP("debug", "d", false, "debug mode")
}

// Get definition using C function call
func definition(word string) string {
	cWord := C.CString(word)
	defer C.free(unsafe.Pointer(cWord))

	cDefinition := C.getDefinition(cWord)
	if cDefinition == nil {
		return ""
	}
	defer C.free(unsafe.Pointer(cDefinition))

	return C.GoString(cDefinition)
}

// Format dictionary output
func output(def string) *string {
	var (
		out      = def
		strArray []string
	)
	out = strings.ReplaceAll(out, " • ", "\n • ")
	out = strings.ReplaceAll(out, " | ", "\n")
	strArray = strings.Split(out, "\n")
	ezlog.Log().M(&strArray).Out()
	return &out
}
