package main

import (
	"context"
	"errors"
	"log"
	pb "spacemesher/proto"
	"strconv"
	"time"

	"google.golang.org/grpc"
)


var tags Tags

type Tags struct {
	Listen string `long:"listen" default:"0.0.0.0:8081" description:"address for listening"`
	Datadir string `short:"d" long:"data-dir" description:"the post data path"`
}



var postcliInfo = make(map[string]map[string]Task)
var PostData *pb.Plot
var server *grpc.Server

type Task struct{
	File_index string
	Duration time.Time
}

type Server struct{
	pb.UnimplementedSpacemesherServer
}


func (s *Server) GetPlot(ctx context.Context, postcli *pb.Postcli) (*pb.Plot, error) {
	file_index, _ := strconv.Atoi(PostData.FileIndex)
	if file_index == int(PostData.TotalFile) { 
		server.Stop() 
		return &pb.Plot{}, errors.New("Done")
	}
	if postcliInfo[postcli.Host] == nil {
		postcliInfo[postcli.Host] = make(map[string]Task)
	}

	task := Task{ File_index: PostData.FileIndex, Duration: time.Now()}
	postcliInfo[postcli.Host][postcli.GPUIndex] = task

	log.Printf("Host: %s GPU: %s\t FileIndex: %s", postcli.Host, postcli.GPUIndex, PostData.FileIndex)
	plot := &pb.Plot{
		NodeId: PostData.NodeId,
		CommitmentAtxId: PostData.CommitmentAtxId,
		DataDir: PostData.DataDir,
		FileIndex: PostData.FileIndex,
		NumUnits: PostData.NumUnits,
		MaxFileSize: PostData.MaxFileSize,
		LabelsPerUnit: PostData.LabelsPerUnit,
		TotalFile: PostData.TotalFile,
	}
	file_index++
	PostData.FileIndex = strconv.Itoa(file_index)
	return plot, nil
}


func (s *Server) GetWorkers(ctx context.Context, _ *pb.StatusReq) (*pb.Workers, error) {
	var workers pb.Workers
	
	for host, v := range postcliInfo{
		var worker pb.Workers_Worker
		worker.Host = host
		for gpu_index, task := range v{
			var provider pb.Workers_Worker_Providers
			provider.Provider = gpu_index
			provider.File_Index = task.File_index
			provider.Duration = time.Now().Sub(task.Duration).String()
			worker.Providers = append(worker.Providers, &provider)
		}
		workers.Worker = append(workers.Worker, &worker)
	}
	
	return &workers, nil
}


func (s *Server) CurrentTask(ctx context.Context, _ *pb.StatusReq) (*pb.Plot, error){
	return PostData, nil
}


func (s *Server) Jump(ctx context.Context, jump *pb.Jump2File) (*pb.StatusResp, error){
	PostData.FileIndex = strconv.Itoa(int(jump.File_Index))
	return &pb.StatusResp{ Success: PostData.FileIndex == strconv.Itoa(int(jump.File_Index)) }, nil
}

