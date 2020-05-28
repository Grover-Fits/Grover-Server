// generate relevant info inside of fits files
// saves images, tables, and header info
// https://fits.gsfc.nasa.gov/fits_primer.html
package main

import (
	"context"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/attron/grover-server/api"
	"github.com/golang/glog"
	"github.com/gorilla/handlers"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joho/godotenv"
	"github.com/siravan/fits"
	"google.golang.org/grpc"
)

// Metadata object
type Metadata struct {
	Filename string
	Metas    []string
	Time     int64
	Images   []string
	Table    string
	Array    string
}

var meta Metadata

// ClientPath -- used to find the locaiton to store images, video, and metadata from accessable location for client
var ClientPath = readEnvFile("CLIENT_PATH")

func readEnvFile(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func readFits(r io.Reader, fn string) Metadata {
	var name string
	var t string
	meta.Filename = fn
	name = meta.Filename
	u, e := fits.Open(r)
	if e != nil {
		log.Fatal(e)
	}

	for i, h := range u {
		//set filename
		out := fmt.Sprintf("/images/%s_%d", name, i)
		fullOut := fmt.Sprintf("%s%s", ClientPath, out)
		// deciding how to handle each HDU
		// XTENSION=IMAGE || SIMPLE=true <- process as image
		// XTENSION=TABLE || XTENSION=BINTABLE <- process as table
		if h.HasImage() {
			t = "image"
			saveImage(h, out)
		} else if h.HasTable() {
			t = "table"
			// 1D array vs nD array
			if len(h.Naxis) == 1 {
				saveArray(h, fullOut)
			} else {
				saveTable(h, fullOut)
			}
		} else {
			t = "Unknown"
			fmt.Println("Unsupported Header Data Unit")
		}
		iStr := strconv.Itoa(i)
		meta.Metas = append(meta.Metas, "/images/meta/"+name+"_"+iStr+".txt")
		metaF, _ := os.OpenFile(ClientPath+meta.Metas[i], os.O_CREATE|os.O_WRONLY, 0644)
		defer metaF.Close()
		fmt.Fprintf(metaF, "***************************************\n")
		fmt.Fprintf(metaF, "HEADER:\t%d\tTYPE:%s\n", i, t)
		fmt.Fprintf(metaF, "***************************************\n")

		// saving header metadata
		for key, value := range h.Keys {
			fmt.Fprintf(metaF, "%s: %v\n", key, value)
		}
	}

	return meta
}

// these functions taken from https://github.com/siravan/fits/blob/master/demo/extract.go
func saveArray(h *fits.Unit, name string) {
	g, _ := os.Create(name + ".dat")
	defer g.Close()

	for i := 0; i < h.Naxis[0]; i++ {
		fmt.Fprintln(g, h.FloatAt(i))
	}
}

func saveImage(h *fits.Unit, name string) {
	n := len(h.Naxis)
	maxis := make([]int, n)
	img := image.NewGray16(image.Rect(0, 0, h.Naxis[0], h.Naxis[1]))
	prod := 1
	for k := 2; k < n; k++ {
		prod *= h.Naxis[k]
	}
	min, max := h.Stats()

	for i := 0; i < prod; i++ {
		l := i
		s := name
		for k := 2; k < n; k++ {
			maxis[k] = l % h.Naxis[k]
			l = l / h.Naxis[k]
			s += fmt.Sprintf("-%d", maxis[k])
		}

		for x := 0; x < h.Naxis[0]; x++ {
			for y := 0; y < h.Naxis[1]; y++ {
				maxis[0] = x
				maxis[1] = y
				if !h.Blank(maxis...) {
					v := uint16((h.FloatAt(maxis...) - min) / (max - min) * 65535) // normalizes based on min and max in the whole image cube
					img.SetGray16(x, h.Naxis[1]-y, color.Gray16{v})
				} else {
					img.SetGray16(x, h.Naxis[1]-y, color.Gray16{0}) // blank pixel
				}
			}
		}

		g, _ := os.Create(ClientPath + s + ".png")
		defer g.Close()
		png.Encode(g, img)
		meta.Images = append(meta.Images, s+".png")
	}
}

func saveTable(h *fits.Unit, name string) {
	g, _ := os.Create(name + ".tab")
	defer g.Close()
	ncols := h.Keys["TFIELDS"].(int)

	label := "" // label is the list of field names/labels
	for col := 0; col < ncols; col++ {
		ttype := h.Keys[fits.Nth("TTYPE", col+1)].(string)
		w := len(h.Format(col, 0)) // the label for each field is resized based on the size of the data on the first row of data
		label += fmt.Sprintf("%-*.*s", w, w, ttype)
	}
	fmt.Fprintln(g, label)

	for row := 0; row < h.Naxis[1]; row++ {
		s := ""
		for col := 0; col < ncols; col++ {
			s += h.Format(col, row)
		}
		fmt.Fprintln(g, s)
	}
}

func run() error {
	go startGrpc()
	grpcServerEndpoint := flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := api.RegisterFitsServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
	if err != nil {
		return err
	}
	handler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{http.MethodPost, http.MethodGet, http.MethodPut, http.MethodDelete}),
		handlers.AllowedHeaders([]string{"Authorization", "Content-Type", "Accept-Encoding", "Accept", "Access-Control-Allow-Origin"}),
	)(mux)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", handler)
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if err := run(); err != nil {
		glog.Fatal(err)
	}
}
