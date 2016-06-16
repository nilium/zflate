// Command zflate is a simple zlib inflate/deflate tool to skip gzip-specific behavior. Useful for
// handling/debugging message compression without headers.
package main

import (
	"compress/zlib"
	"flag"
	"io"
	"log"
	"os"
)

var (
	infile  = flag.String("i", "-", "Input `file`. Defaults to standard input ('-').")
	outfile = flag.String("o", "-", "Output `file`. Defaults to standard output ('-').")
	level   = flag.Int("l", zlib.DefaultCompression, "The compression `level` to use when deflating.")
	inflate = flag.Bool("d", false, "Inflate input.")
)

func main() {
	log.SetPrefix("zflate: ")
	log.SetFlags(0)

	flag.Parse()
	if flag.NArg() > 0 {
		log.Print("invalid arguments")
		flag.Usage()
		os.Exit(0)
	}

	var in io.ReadCloser = os.Stdin
	if *infile != "-" {
		fi, err := os.Open(*infile)
		if err != nil {
			log.Panic("Unable to open input file ", *infile, ": ", err)
		}
		defer fi.Close()
		in = fi
	}

	var out io.WriteCloser = os.Stdout
	if *outfile != "-" {
		fi, err := os.Create(*outfile)
		if err != nil {
			log.Panic("Unable to open output file ", *outfile, ": ", err)
		}
		defer fi.Close()
		out = fi
	}

	var err error

	op := "inflate"
	if *inflate {
		in, err = zlib.NewReader(in)
		if err != nil {
			log.Panic("Unable to create zlib reader for input: ", err)
		}
		defer in.Close()
	} else {
		op = "deflate"
		out = zlib.NewWriter(out)
		defer out.Close()
	}

	if _, err = io.Copy(out, in); err != nil {
		log.Panic("Unable to ", op, " input: ", err)
	}
}
