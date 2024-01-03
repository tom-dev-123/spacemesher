package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	pb "spacemesher/proto"
)

func Decode(code string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return "", err
	}
	decode_str := hex.EncodeToString(bytes)
	return decode_str, nil
}

func NewPlot(post_data_dir string) error {
	data, err := ioutil.ReadFile(post_data_dir + "/postdata_metadata.json")
	if err != nil {
		return err
	}
	var postdata pb.Plot
	err = json.Unmarshal(data, &postdata)
	if err != nil{
		return err
	}
	atx, _ := Decode(postdata.CommitmentAtxId)
	nodeid, _ := Decode(postdata.NodeId)
	PostData = &pb.Plot{
		NodeId: nodeid, 
		CommitmentAtxId: atx, 
		DataDir: post_data_dir,
		FileIndex: "0",
		NumUnits: postdata.NumUnits,
		LabelsPerUnit: postdata.LabelsPerUnit,
		MaxFileSize: postdata.MaxFileSize,
		TotalFile: postdata.NumUnits * 16,
	}
	log.Println(PostData)
	return nil
}
