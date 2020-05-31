package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/attron/grover-server/api"
	"google.golang.org/grpc"
)

type server struct{}

// GetMovie -- send request to generate movie from images
func (*server) GetMovie(ctx context.Context, req *api.GetMovieRequest) (*api.GetMovieResponse, error) {
	log.Println("Received GetMovie request!")
	files := req.GetFilePath()
	// need to break out files and append client path so command can locate the images
	sepFiles := strings.Fields(files)
	sendFiles := ""
	for _, s := range sepFiles {
		sendFiles = sendFiles + " " + ClientPath + s
	}
	out, err := exec.Command("./convert-to-video.sh", sendFiles, ClientPath).Output()
	if err != nil {
		log.Println(err)
	}
	log.Println("Received files for movie generation: " + sendFiles)
	return &api.GetMovieResponse{
		MovLoc: string(out),
	}, nil
}

// UploadFitsFiles -- called when new fits files have been uploaded
func (*server) UploadFitsFiles(ctx context.Context, req *api.UploadFitsFilesRequest) (*api.UploadFitsFilesResponse, error) {
	log.Println("Received UploadFile request!")
	file := req.GetFileContent()
	name := req.GetName()
	sDec, err := base64.StdEncoding.DecodeString(file)
	if err != nil {
		fmt.Printf("Error decoding string: %s ", err.Error())
		return nil, err
	}

	s := strings.NewReader(string(sDec))
	fmt.Println("attempting to send fits data . . .", name)
	meta := handleIncoming(s, name)
	out, err := json.Marshal(meta)
	if err != nil {
		log.Fatalf("Failed to retireve metadata from file!")
	}
	meta = Metadata{}
	return &api.UploadFitsFilesResponse{
		Metadata: string(out),
	}, nil
}

func (*server) GetFitsFiles(ctx context.Context, req *api.GetFitsFilesRequest) (*api.GetFitsFilesResponse, error) {
	return nil, nil
}

func (*server) TestClient(ctx context.Context, req *api.TestClientRequest) (*api.TestClientResponse, error) {
	fmt.Println("Received TestClient request!")
	msg := req.GetMsg()

	return &api.TestClientResponse{
		Msg: "Your message was: " + msg,
	}, nil
}

func startGrpc() {
	fmt.Println("Attempting to start GRPC server . . .")
	lis, err := net.Listen("tcp", "127.0.0.1:9090")
	if err != nil {
		log.Fatalf("failed to start listener: %v", err)
	}
	fmt.Println("Server started successfully!")
	opts := []grpc.ServerOption{grpc.MaxRecvMsgSize(104857600)}
	gServ := grpc.NewServer(opts...)
	api.RegisterFitsServiceServer(gServ, &server{})
	gServ.Serve(lis)
}
