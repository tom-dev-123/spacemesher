package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	pb "spacemesher/proto"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

type Client struct{
	conn *grpc.ClientConn
	client pb.SpacemesherClient
	ip string
	close chan struct{}
}

type Task struct {
	Command *exec.Cmd
}


var tag Tags

type Tags struct {
	Address string `long:"address" default:"0.0.0.0:8081" description:"Address for listening"`
	Provider uint `short:"p" long:"provider" description:"Binging gpu"`
	Postcli string `long:"bin" default:"postcli" description:"the path of postcli binary"`
}

func New() *Client {
	return &Client{
		ip: GetLocalIP(),
		close: make(chan struct{}),
	}
}


func (c *Client)Connect(addr string) error {
	var conn *grpc.ClientConn
	for {
		conn, _ = grpc.Dial(addr, grpc.WithInsecure())
		conn.WaitForStateChange(context.Background(), conn.GetState())
		switch conn.GetState(){
		case connectivity.Connecting:
			log.Print("Connecting...")
			time.Sleep(time.Second * 5)
			continue
		case connectivity.Ready:
			c.client = pb.NewSpacemesherClient(conn)
			c.conn = conn
			return nil
		case connectivity.TransientFailure:
			log.Println("Waiting server up...")
			time.Sleep(time.Second * 5)
			continue
		case connectivity.Idle:
			log.Println("Client idle, please restart!")
			c.Close()
			return nil
		default:
			log.Println(conn.GetState())
			time.Sleep(5 * time.Second)
			continue
		}
	}
}

func (c *Client) Task_Start() {
	for{
		if c.conn.GetState() != connectivity.Ready { c.Connect(tag.Address) }
		plot, err := c.Get()
		if err != nil {
			switch err.Error(){
			case "Done":
				continue
			default:
				log.Printf("error: %s", err.Error())
				continue
			}
		}
		postcli, err := exec.LookPath(tag.Postcli)
		if err != nil {
			log.Fatal(err)
		}
		_, err = os.Stat(plot.DataDir)
		for os.IsNotExist(err) { log.Print(err); time.Sleep(time.Second * 5)}
		
		task := Task{
			Command: exec.Command(
				postcli,
				"-id",
				plot.NodeId,
				"-commitmentAtxId", 
				plot.CommitmentAtxId, 
				"-datadir",
				plot.DataDir,
				"-labelsPerUnit",
				strconv.Itoa(int(plot.LabelsPerUnit)),
				"-maxFileSize",
				strconv.Itoa(int(plot.MaxFileSize)),
				"-numUnits",
				strconv.Itoa(int(plot.NumUnits)),
				"-fromFile",
				plot.FileIndex,
				"-toFile",
				plot.FileIndex,
				"-provider",
				strconv.Itoa(int(tag.Provider)),
			),
		}
		log.Println(task.Command.Args)
		
		// _, err = c.Plotting(&plotting)
		// if err != nil {
		// 	log.Fatal(err)
		// }
		err = task.Command.Start()
		task.Command.Wait()
	}
}

func (c *Client) Get() (*pb.Plot, error) {
	plot := new(pb.Plot)
	plotting := pb.Postcli{
		Host: c.ip,
		GPUIndex: strconv.Itoa(int(tag.Provider)),
	}

	plot, err := c.client.GetPlot(context.Background(), &plotting)

	if err != nil {
		return &pb.Plot{}, err
	}
	return plot, nil
}

// func (c *Client) Plotting(postcli *pb.Postcli) (bool, error) {
// 	resp, err := c.client.Plotting(context.Background(), postcli)
// 	return resp.Success, err
// }

func (c *Client) Close(){
	c.close <- struct{}{}
}

